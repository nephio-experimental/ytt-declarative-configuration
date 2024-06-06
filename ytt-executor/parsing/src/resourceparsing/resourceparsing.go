/*
Copyright 2022. Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
