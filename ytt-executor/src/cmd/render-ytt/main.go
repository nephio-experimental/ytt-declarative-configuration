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

package main

import (
	"fmt"
	"os"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/commandExec"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/process"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	yttProc := YttProcessor{}
	cmd := command.Build(&yttProc, command.StandaloneEnabled, false)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type YttProcessor struct{}

func (yttProc *YttProcessor) Process(resourceList *framework.ResourceList) error {
	// Check for config
	if resourceList.FunctionConfig.IsNilOrEmpty() {
		logger.LogWarning("No function config provided. Default values will be used.")
	} else {
		// Get and parse config
		err := config.Configure(resourceList.FunctionConfig)
		if err != nil {
			resourceList.Results = logger.LogStack
			return err
		}
	}

	// Write kpt input to file system
	fileArgs, baseDir, err := process.ParseAndWriteKYamlRNodesAsYttTemplates(resourceList.Items...)
	if err != nil {
		resourceList.Results = logger.LogStack
		return err
	}

	//Remove temporary directory of ytt files.
	defer os.RemoveAll(baseDir)

	// Execute ytt binary with given file arguments
	yttOutputBuffer, err := commandExec.ExecuteYttForTemplate(fileArgs)
	if err != nil {
		resourceList.Results = logger.LogStack
		return err
	}

	// Collect output RNodes
	// TODO: if expanding / changing move to process package
	var outputItems []*kyaml.RNode
	for _, item := range resourceList.Items {
		if item.GetKind() == config.YttOutputFileKind {
			outputItems = append(outputItems, item)
		}
	}

	// Take ytt executable output and parse back to kyaml.RNode
	err = process.UnmarshalYttOutput(yttOutputBuffer, outputItems)
	if err != nil {
		resourceList.Results = logger.LogStack
		return err
	}

	resourceList.Results = logger.LogStack
	return nil
}
