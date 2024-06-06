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

package variablemapping

import (
	"fmt"
	"parsing/types"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	ktypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

func ParseSubstitutions(variableMapping *fn.KubeObject) (types.Substitutions, error) {
	subs := types.Substitutions{}
	subs.Substitutions = []*types.Substitution{}
	err := variableMapping.As(&subs)

	if err != nil {
		return subs, err
	}

	return subs, nil
}

// Do as Generic ParseFileId<T> T = ConfigFileId, OutputFileId
// OR take field name as input (configFile, outputFile and use NestedResource)
func ParseConfigFile(variableMapping *fn.KubeObject) (types.ConfigFile, error) {
	configFile := types.ConfigFile{}
	err := variableMapping.As(&configFile)

	if err != nil {
		return configFile, err
	}

	return configFile, nil
}

func ParseOutputFile(variableMapping *fn.KubeObject) (types.OutputFile, error) {
	configFile := types.OutputFile{}
	err := variableMapping.As(&configFile)

	if err != nil {
		return configFile, err
	}

	return configFile, nil
}

func NewSource() ktypes.SourceSelector {
	return ktypes.SourceSelector{
		Options: &ktypes.FieldOptions{},
	}
}

func ParseSubConfigs(
	substitutions types.Substitutions,
	resourceMap map[resid.ResId]*fn.KubeObject) ([]types.SourceMeta, error) {
	sources := []types.SourceMeta{}

	for _, sub := range substitutions.Substitutions {
		subResource := resourceMap[sub.Replacement]
		source := NewSource()
		ok, err := subResource.NestedResource(&source, "data")
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("Failed to read 'data' field from %v", sub.Replacement)
		}
		source.ResId = sub.Replacement
		sources = append(sources, types.SourceMeta{
			Substitution: sub,
			Source:       &source,
		})
	}

	return sources, nil
}

func UpdateApplyReplacementsFnConfig(
	applyReplacements *fn.KubeObject, replacements []ktypes.Replacement,
) error {
	if len(replacements) == 0 {
		_, err := applyReplacements.RemoveNestedField("replacements")
		return err
	}
	return applyReplacements.SetNestedField(replacements, "replacements")
}

func ParseTemplateConfig(templateConfig *fn.KubeObject) (map[string]int, error) {

	m := make(map[string]int, 0)
	err := templateConfig.GetMap("data").As(&m)

	if err != nil {
		return m, err
	}
	return m, nil
}
