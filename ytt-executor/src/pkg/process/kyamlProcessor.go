// Copyright 2024 Ericsson AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/fileWriter"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

// templateType internal enum for handling file processing type
// To distinguish if file should use -f or --data-values-file argument
//
// defaultTemplate: For most  basic template processing, file name is returned with -f argument
//
// valuesTemplate: For explicitly specifying values file, file name is returned with --data-values-file argument
type templateType int

const (
	defaultTemplate templateType = iota
	valuesTemplate
	outputFile
)

// ParseAndWriteKYamlRNodesAsYttTemplates
//
// Parameters:
//   - items: list of yaml.RNode items to write to work directory for ytt binary processing
//
// Returns:
//   - fileArgs: Any error that could be experienced when writing the file
//   - error: Any error that could be experienced when writing the file
func ParseAndWriteKYamlRNodesAsYttTemplates(items ...*kyaml.RNode) (fileArgs []string, baseDir string, err error) {
	baseDir, err = os.MkdirTemp("", "ytt-files-*")
	if err != nil {
		return []string{}, "", fmt.Errorf("Directory creation for ytt files failed: %v", err)
	}
	for _, item := range items {

		// Check for file type
		itemType := getItemTemplateType(item)

		// Switch based on file type
		switch itemType {

		// Write file and return -f <file_name> argument
		case defaultTemplate:
			fileName, err := processKYamlRNode(item, baseDir)
			if err != nil {
				return fileArgs, baseDir, err
			}
			fileArgs = append(fileArgs, "-f", fileName)
			break

		// Write file and return --data-values-file <file_name> argument
		case valuesTemplate:
			fileName, err := processKYamlRNode(item, baseDir)
			if err != nil {
				return fileArgs, baseDir, err
			}
			fileArgs = append(fileArgs, "--data-values-file", fileName)
			break

		//	Output file should not be written or handled by ytt bin
		case outputFile:
			break
		}
	}
	return fileArgs, baseDir, nil
}

// getItemTemplateType function to identify template type based on configuration
//
// Parameters:
//   - item: yaml.RNode to be identified
//
// Returns:
//   - templateType: enum identifying template type
func getItemTemplateType(item *kyaml.RNode) templateType {
	// Default check for values File by Kind
	if config.YttInputValuesFileHandling == config.ValuesIdentifierKind {
		if item.GetKind() == config.YttInputValueFileKind {
			return valuesTemplate
		}
	}

	if item.GetKind() == config.YttOutputFileKind {
		return outputFile
	}
	return defaultTemplate
}

func processKYamlRNode(item *kyaml.RNode, baseDir string) (fileName string, err error) {
	fileName = baseDir + "/" + item.GetAnnotations()["config.kubernetes.io/path"]

	// Log detailed info about files
	logger.LogDetailedDebug(fmt.Sprintf("Writing file for ytt processing: %s", fileName), map[string]string{
		"kyaml":         item.MustString(),
		"fileName":      fileName,
		"annotationKey": config.YttNodeAnnotations,
		"hasAnnotation": strconv.FormatBool(!item.Field(config.YttNodeAnnotations).IsNilOrEmpty()),
		"contentKey":    config.YttNodeContent,
		"hasData":       strconv.FormatBool(!item.Field(config.YttNodeContent).IsNilOrEmpty()),
	})

	// Hande annotations
	if !item.Field(config.YttNodeAnnotations).IsNilOrEmpty() {

		// Get string of annotation data
		yttAnnotationData := item.Field(config.YttNodeAnnotations).Value.MustString()

		// strip first line (no other way about it)
		yttAnnotationData = yttAnnotationData[strings.Index(yttAnnotationData, "\n")+1:]
		err := fileWriter.WriteToFile(fileName, yttAnnotationData)
		if err != nil {
			return fileName, err
		}
	}

	// Handle content
	if !item.Field(config.YttNodeContent).IsNilOrEmpty() {
		err := fileWriter.WriteToFile(fileName, item.Field(config.YttNodeContent).Value.MustString())
		if err != nil {
			return fileName, err
		}
	}

	return fileName, nil
}
