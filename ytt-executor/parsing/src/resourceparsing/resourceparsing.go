package resourceparsing

import (
	"encoding/json"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func NewResIdWithNameAndKind(name string, kind string) yaml.ResourceIdentifier {
	return yaml.ResourceIdentifier{
		TypeMeta: yaml.TypeMeta{
			Kind: kind,
		},
		NameMeta: yaml.NameMeta{
			Name: name,
		},
	}
}

func CreateFileMap(resourceList *framework.ResourceList) (map[yaml.ResourceIdentifier]*yaml.RNode, error) {
	fileMap := map[yaml.ResourceIdentifier]*yaml.RNode{}

	for _, item := range resourceList.Items {
		fileMeta, err := item.GetMeta()
		fileId := fileMeta.GetIdentifier()
		if err != nil {
			return nil, err
		}
		fileMap[NewResIdWithNameAndKind(fileId.Name, fileId.Kind)] = item
	}

	return fileMap, nil
}

func CreateFileMapIndex(resourceList *framework.ResourceList) (map[yaml.ResourceIdentifier]int, error) {
	fileMap := map[yaml.ResourceIdentifier]int{}

	for i, item := range resourceList.Items {
		fileMeta, err := item.GetMeta()
		fileId := fileMeta.GetIdentifier()
		if err != nil {
			return nil, err
		}
		fileMap[NewResIdWithNameAndKind(fileId.Name, fileId.Kind)] = i
	}

	return fileMap, nil
}

type IdWrapper interface {
	GetName() string
	GetKind() string
}

type IdWrapperImpl struct {
	Id ResourceIdentifier
}

func (w IdWrapperImpl) GetName() string {
	return w.Id.Name
}

func (w IdWrapperImpl) GetKind() string {
	return w.Id.Kind
}

type OutputFile struct {
	Id ResourceIdentifier `json:"outputFile" yaml:"outputFile"`
}

type NDCIQFile struct {
	Id ResourceIdentifier `json:"networkDesignerCIQ" yaml:"networkDesignerCIQ"`
}

type ITCIQFile struct {
	Id ResourceIdentifier `json:"integrationTechnicianCIQ" yaml:"integrationTechnicianCIQ"`
}

type ResourceIdentifier struct {
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
}

func ParseGeneratedConfig(generatedConfig *yaml.RNode) (yaml.ResourceIdentifier, error) {

	bytes, err := generatedConfig.MarshalJSON()
	if err != nil {
		return yaml.ResourceIdentifier{}, err
	}
	of := OutputFile{}
	err = json.Unmarshal(bytes, &of)

	return NewResIdWithNameAndKind(of.Id.Name, of.Id.Kind), nil
}

func ParseNDCIQFile(generatedConfig *yaml.RNode) (yaml.ResourceIdentifier, error) {

	bytes, err := generatedConfig.MarshalJSON()
	if err != nil {
		return yaml.ResourceIdentifier{}, err
	}
	of := NDCIQFile{}
	err = json.Unmarshal(bytes, &of)

	return NewResIdWithNameAndKind(of.Id.Name, of.Id.Kind), nil
}

func ParseITCIQFile(generatedConfig *yaml.RNode) (yaml.ResourceIdentifier, error) {

	bytes, err := generatedConfig.MarshalJSON()
	if err != nil {
		return yaml.ResourceIdentifier{}, err
	}
	of := ITCIQFile{}
	err = json.Unmarshal(bytes, &of)

	return NewResIdWithNameAndKind(of.Id.Name, of.Id.Kind), nil
}
