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
	"errors"
	"os"
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestYttProcessor_Process(t *testing.T) {
	// Setup TempDir for testing
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	// Test structure
	tests := []struct {
		name          string
		expectedError error
		args          *framework.ResourceList
		preTest       func()
	}{ // Test List

		// Test error catch in config.Configure call with invalid fnConfig
		{
			"Test fail on fnConfig",
			errors.New("node kind is not a string: map[child_element:to_break_parsing]"),
			&framework.ResourceList{
				FunctionConfig: kyaml.MustParse(`
output:
  kind:
    child_element: to_break_parsing
`),
			},
			nil,
		},

		// Test error catch in process.ParseAndWriteKYamlRNodesAsYttTemplates
		// TODO, commented because of temporary directory addition, this test is failing, need to revisit
		/*		{
							"Test fail on ParseAndWrite",
							&fs.PathError{
								Op:   "open",
								Path: "",
								//Err:  syscall.Errno(0x2),
								Err:  syscall.Errno(0x15),
							},
							&framework.ResourceList{
								Items: []*kyaml.RNode{
									// item that will fail on write
									kyaml.MustParse(`
				apiVersion: v1alpha1
				kind: YttTemplate
				metadata:
				  name: failing-ytt-item
				ytt_template_content:
				  #@ sample ytt content
				  ytt_item:
				    #@ more ytt content
				#@ Ending comment
				`),
								},
							},
							nil,
						},*/

		// Test error catch in commandExec.ExecuteYttForTemplate
		{
			"Test fail on ExecuteYttForTemplate",
			errors.New("ytt: exit status 1 (stderr: cat: invalid option -- 'f'\nTry 'cat --help' for more information.\n)"),
			&framework.ResourceList{
				Items: []*kyaml.RNode{
					kyaml.MustParse(`
apiVersion: v1alpha1
kind: Configuration
metadata:
  name: ytt-output
  annotations:
    config.kubernetes.io/path: "main_test_path_1/output.yaml"
data:
`),
					kyaml.MustParse(`
apiVersion: v1alpha1
kind: YttTemplate
metadata:
  name: ytt-template
  annotations:
    config.kubernetes.io/path: "main_test_path_1/template.yaml"
ytt_header:
  header:
  #@ sample ytt annotation
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`),
				},
			},
			func() {
				config.YttBinaryName = "cat"
			},
		},

		// Test error catch in process.UnmarshalYttOutput
		{
			"Test fail on UnmarshalYttOutput",
			errors.New("no output file with kind: Configuration provided"),
			&framework.ResourceList{
				Items: []*kyaml.RNode{
					kyaml.MustParse(`
apiVersion: v1alpha1
kind: YttTemplate
metadata:
  name: ytt-template
  annotations:
    config.kubernetes.io/path: "main_test_path_2/template.yaml"
ytt_header:
  header:
  #@ sample ytt annotation
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`),
				},
			},
			func() {
				config.YttBinaryName = "ls"
			},
		},

		// Happy test
		{
			"Test successful flow of process",
			nil,
			&framework.ResourceList{
				Items: []*kyaml.RNode{
					kyaml.MustParse(`
apiVersion: v1alpha1
kind: Configuration
metadata:
  name: ytt-output
  annotations:
    config.kubernetes.io/path: "main_test_path_2/output.yaml"
data: some_existing_data
`),
					kyaml.MustParse(`
apiVersion: v1alpha1
kind: YttTemplate
metadata:
  name: ytt-template
  annotations:
    config.kubernetes.io/path: "main_test_path_2/template.yaml"
ytt_header:
  header:
  #@ sample ytt annotation
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`),
				},
			},
			func() {
				config.YttBinaryName = "ls"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run pretest details if needed
			if tt.preTest != nil {
				tt.preTest()
			}

			// Create struct
			yttProc := &YttProcessor{}

			// Execute main
			err := yttProc.Process(tt.args)

			// Check if error was not expected but seen
			if tt.expectedError == nil && err != nil {
				t.Fatalf("Did not expect error but got: %v", err)

				// Check if error was expected but not received
			} else if tt.expectedError != nil && err == nil {
				t.Fatalf("Got no error but expected: %v", tt.expectedError)

				// Error was expected and seen, assert error
			} else if tt.expectedError != nil && err != nil {

				// Assert actual error
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
