package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	goexec "os/exec"
	rparsing "parsing/resourceparsing"
	"strings"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

type OutputKeyMeta struct {
	// Name is the metadata.name field of a Resource
	Key string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ConfigResource struct {
	TypeMeta kyaml.TypeMeta         `json:",inline" yaml:",inline"`
	MetaData kyaml.NameMeta         `json:"metadata" yaml:"metadata"`
	Config   map[string]interface{} `json:"values" yaml:"values"`
}

type FnConfig struct {
	Template *kyaml.ResourceIdentifier   `json:"template" yaml:"template"`
	Output   *kyaml.ResourceIdentifier   `json:"output" yaml:"output"`
	Schemas  []*kyaml.ResourceIdentifier `json:"schemas" yaml:"schemas"`
	Ciqs     []*kyaml.ResourceIdentifier `json:"ciqs" yaml:"ciqs"`
}

func newConfigResource(kind string, name string) ConfigResource {

	return ConfigResource{
		TypeMeta: kyaml.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       kind,
		},
		MetaData: kyaml.NameMeta{
			Name: name,
		},
		Config: make(map[string]interface{}, 0),
	}
}

func main() {
	// return // TMP nil // TMP

	gsp := GenerateSettersProcessor{}
	cmd := command.Build(&gsp, command.StandaloneEnabled, false)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type GenerateSettersProcessor struct{}

func trimQuotes(s string) string {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		print("trimQuotes")
		// Remove the first and last characters
		return s[1 : len(s)-1]
	}
	if strings.HasPrefix(s, " |") {
		print("trimSpacePipe")
		// Remove the first two characterss
		return s[2:len(s)]
	}
	return s
}

func createYttString(node *kyaml.RNode, field string, prependStr string, postpendStr string) (string, error) {
	fieldMatcher := kyaml.Get(field)
	node, err := fieldMatcher.Filter(node)
	if err != nil {
		return "", err
	}
	nodeStr, err := node.String()
	if err != nil {
		return "", err
	}
	nodeStr = trimQuotes(nodeStr)
	return fmt.Sprintf("%s%s%s", prependStr, nodeStr, postpendStr), nil
}

func WriteToFile(strFile string, category string) (string, error) {
	f, err := os.CreateTemp("", "ytt-"+category+"-*.yaml")
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(strFile)
	return f.Name(), err
}

type contentCategory struct {
	content  string
	category string
}

func createYttTmpFiles(yttStrFiles []contentCategory) ([]string, error) {
	fileNames := make([]string, len(yttStrFiles))
	for i, strFile := range yttStrFiles {

		fileName, err := WriteToFile(strFile.content, strFile.category)
		if err != nil {
			return fileNames, err
		}

		fileNames[i] = fileName
	}
	return fileNames, nil
}

func removeFiles(fileNames []string) {
	for _, fileName := range fileNames {
		os.Remove(fileName)
	}
}

func process(resourceList *framework.ResourceList) error {
	conf, err := parseFnConfig(resourceList.FunctionConfig)
	if err != nil {
		return err
	}
	rMap, err := rparsing.CreateFileMapIndex(resourceList)

	if err != nil {
		return err
	}

	var fileStrs []contentCategory
	var schemaStrs []contentCategory

	for i, _ := range conf.Schemas {
		schemaStr, err := createYttString(resourceList.Items[rMap[*conf.Schemas[i]]], "schema", "\n#@ load(\"@ytt:overlay\", \"overlay\")\n#@data/values-schema\n---\n#@overlay/match missing_ok=True\n", "")
		if err != nil {
			return err
		}
		fileStrs = append(fileStrs, contentCategory{schemaStr, "schema"})
		schemaStrs = append(schemaStrs, contentCategory{schemaStr, "schema"})
	}

	ciqStr, err := createYttString(resourceList.Items[rMap[*conf.Ciqs[0]]], "values", "#@data/values\n---\n", "")
	if err != nil {
		return err
	}
	fileStrs = append(fileStrs, contentCategory{ciqStr, "ciq"})

	templateStr, err := createYttString(resourceList.Items[rMap[*conf.Template]], "template", "#@ load(\"@ytt:data\", \"data\")\n", "")
	if err != nil {
		return err
	}
	fileStrs = append(fileStrs, contentCategory{templateStr, "template"})

	var fileInput = ""

	for _, fileStr := range fileStrs {
		fileInput = fileInput + "\n" + fileStr.category + ":\n==========\n" + fileStr.content
	}

	fileNames, err := createYttTmpFiles(fileStrs)
	defer removeFiles(fileNames)
	if err != nil {
		return err
	}

	ytt := &Ytt{}
	outstr, err := ytt.Execute(fileInput, fileNames)

	if err != nil {
		return err
	}

	configId := rparsing.NewResIdWithNameAndKind(conf.Output.Name, conf.Output.Kind)
	cr := newConfigResource(configId.Kind, configId.Name)
	err = yaml.Unmarshal([]byte(outstr), cr.Config)

	yamlbytes, err := yaml.Marshal(cr)
	if err != nil {
		return err
	}

	rnode, err := kyaml.Parse(string(yamlbytes))
	configIndex, oldConfigExists := rMap[configId]

	if oldConfigExists {
		resourceList.Items[configIndex] = rnode
	} else {
		print("new ciq config")
		resourceList.Items = append(resourceList.Items, rnode)
	}

	// OpenAPI

	fileNames2, err := createYttTmpFiles(schemaStrs)
	defer removeFiles(fileNames2)
	if err != nil {
		return err
	}

	outstr, err = ytt.RenderOpenAPI(fileNames2)

	if err != nil {
		return err
	}

	openId := rparsing.NewResIdWithNameAndKind("open-api-schema", "OpenAPISchema")
	openR := newConfigResource(openId.Kind, openId.Name)
	err = yaml.Unmarshal([]byte(outstr), openR.Config)
	yamlbytes, err = yaml.Marshal(openR)
	if err != nil {
		return err
	}

	rnode, err = kyaml.Parse(string(yamlbytes))
	openIndex, oldOpenExists := rMap[openId]

	if oldOpenExists {
		resourceList.Items[openIndex] = rnode
	} else {
		print("new OpenAPI config")
		resourceList.Items = append(resourceList.Items, rnode)
	}

	return nil
}

func parseFnConfig(fnConfigNode *kyaml.RNode) (*FnConfig, error) {
	bytes, err := fnConfigNode.MarshalJSON()
	if err != nil {
		return &FnConfig{}, err
	}
	rI := FnConfig{}
	err = json.Unmarshal(bytes, &rI)

	return &rI, nil
}

func (gsp *GenerateSettersProcessor) Process(resourceList *framework.ResourceList) error {
	return process(resourceList)
}

type Ytt struct{}

func (ytt *Ytt) Execute(fileInput string, files []string) (string, error) {
	args := make([]string, len(files))

	for i, fileName := range files {
		args[i] = fmt.Sprintf("-f%s", fileName)
	}

	var stdoutBs, stderrBs bytes.Buffer
	var stdin io.Reader

	cmd := goexec.Command("ytt", args...)
	cmd.Stdin = stdin
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()
	if err != nil {
		stderrStr := stderrBs.String()
		return "", fmt.Errorf("Executing ytt: %s (stderr: %s - input: %s)", err, stderrStr, fileInput)
	}

	return stdoutBs.String(), nil
}

func (ytt *Ytt) RenderOpenAPI(schemas []string) (string, error) {
	args := make([]string, len(schemas)+2)

	for i, fileName := range schemas {
		args[i] = fmt.Sprintf("-f%s", fileName)
	}
	args[len(args)-2] = "--data-values-schema-inspect"
	args[len(args)-1] = "--output=openapi-v3"

	// fmt.Println(args)

	var stdoutBs, stderrBs bytes.Buffer
	var stdin io.Reader

	cmd := goexec.Command("ytt", args...)
	cmd.Stdin = stdin
	cmd.Stdout = &stdoutBs
	cmd.Stderr = &stderrBs

	err := cmd.Run()
	if err != nil {
		stderrStr := stderrBs.String()
		return "", fmt.Errorf("Executing ytt: %s (stderr: %s)", err, stderrStr)
	}

	return stdoutBs.String(), nil
}
