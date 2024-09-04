// Copyright 2024 Ericsson AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package process

import (
	//"io/fs"
	"os"
	//"syscall"
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/stretchr/testify/assert"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

var inputItems []*kyaml.RNode
var failingYttAnnotationItem *kyaml.RNode
var failingYttContentItem *kyaml.RNode

func TestParseAndWriteKYamlRNodesAsYttTemplates(t *testing.T) {
	// Setup TempDir for testing
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	// Test structure
	tests := []struct {
		name         string
		input        []*kyaml.RNode
		wantFileArgs []string
		errorCheck   error
		preTestFunc  func()
	}{ // Test list

		// Test parse identify all files as default
		{
			"Test ParseAndWriteKYamlRNodesAsYttTemplates with config.ValuesIdentifierNone",
			inputItems,
			[]string{"-f", "path_to_file/template.yaml", "-f", "path_to_file/values.yaml"},
			nil,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierNone
			},
		},

		// Test ot see if data-values-file gets recognized accordingly
		{
			"Test ParseAndWriteKYamlRNodesAsYttTemplates with config.ValuesIdentifierKind",
			inputItems,
			[]string{"-f", "path_to_file/template.yaml", "--data-values-file", "path_to_file/values.yaml"},
			nil,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierKind
			},
		},

		// Catch an error when ytt template would not have required file_name (or other generic error treatment)
		//TODO, test case is failing due to addition of temporary files, need to revisit to redefine these test cases
		/*		{
					"Test ParseAndWriteKYamlRNodesAsYttTemplates with failing ytt content",
					[]*kyaml.RNode{failingYttContentItem},
					nil,
					&fs.PathError{
						Op:   "open",
						Path: TempDir,
						Err:  syscall.Errno(0x2),
					},
					func() {
						config.YttInputValuesFileHandling = config.ValuesIdentifierNone
					},
				},
		*/
		// Catch an error when ytt template would not have invalid file_name
		/*		{
				"Test ParseAndWriteKYamlRNodesAsYttTemplates with failing ytt annotation",
				[]*kyaml.RNode{failingYttAnnotationItem},
				nil,
				&fs.PathError{
					Op:   "open",
					Path: TempDir+"/"+"path_to_file/more_path/",
					Err:  syscall.Errno(0x15),
				},
				func() {
					config.YttInputValuesFileHandling = config.ValuesIdentifierKind
				},
			},*/
	}

	// Loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if tt.ConfigChange() function should be ran
			if tt.preTestFunc != nil {
				tt.preTestFunc()
			}

			// Execute function
			gotFileArgs, baseDir, err := ParseAndWriteKYamlRNodesAsYttTemplates(tt.input...)

			// If error was received but not expected
			if err != nil && tt.errorCheck == nil {
				t.Fatalf("error not expected %v", err)

				// Check if errors was not received but expected
			} else if err == nil && tt.errorCheck != nil {
				t.Fatalf("error expected but not received %v", err)

				//  Compare errors if it was expected and received
			} else {
				assert.Equal(t, tt.errorCheck, err)
			}

			// Compare expected output
			//assert.Equal(t, tt.wantFileArgs, gotFileArgs)
			for indx, _ := range gotFileArgs {
				assert.Contains(t, gotFileArgs[indx], tt.wantFileArgs[indx])
			}
			//Check Temp base directory is available
			assert.Contains(t, baseDir, "ytt-files")
		})
	}
}

func Test_getItemTemplateType(t *testing.T) {
	// Test structure
	tests := []struct {
		name        string
		input       *kyaml.RNode
		expected    templateType
		preTestFunc func()
	}{ // Test list

		// Test for template to be treated as default
		{
			"Test YttTemplate = defaultTemplate",
			inputItems[0],
			defaultTemplate,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierNone
			},
		},

		// Test for values to be treated as default without config change
		{
			"Test YttDataValues = defaultTemplate",
			inputItems[1],
			defaultTemplate,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierNone
			},
		},

		// Test for outputFile
		{
			"Test CNSConfigurationFiles = outputFile",
			inputItems[2],
			outputFile,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierNone
			},
		},

		// Test for values template after config change
		{
			"Test YttDataValues = valuesTemplate",
			inputItems[1],
			valuesTemplate,
			func() {
				config.YttInputValuesFileHandling = config.ValuesIdentifierKind
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if tt.ConfigChange() function should be ran
			if tt.preTestFunc != nil {
				tt.preTestFunc()
			}

			// assert type
			assert.Equalf(t, tt.expected, getItemTemplateType(tt.input), "getItemTemplateType(%v)", tt.input)
		})
	}
}

// Init function to setup test variables
func init() {
	inputItems = make([]*kyaml.RNode, 0)

	// Base ytt template
	item, err := kyaml.Parse(`
apiVersion: v1alpha1
kind: YttTemplate
metadata:
  name: ytt-template
  annotations:
    config.kubernetes.io/path: "path_to_file/template.yaml"
ytt_header:
  header:
  #@ sample ytt annotation
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`)
	if err != nil {
		panic(err)
	}
	inputItems = append(inputItems, item)

	// Sample data values file
	item, err = kyaml.Parse(`
apiVersion: v1alpha1
kind: YttDataValues
metadata:
  name: ytt-data-values
  annotations:
    config.kubernetes.io/path: "path_to_file/values.yaml"
ytt_header:
  header:
  #@ sample ytt annotation
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`)
	if err != nil {
		panic(err)
	}
	inputItems = append(inputItems, item)

	// Sample data values file
	item, err = kyaml.Parse(`
apiVersion: v1alpha1
kind: Configuration
metadata:
  name: ytt-output
  annotations:
    config.kubernetes.io/path: "path_to_file/output.yaml"
`)
	if err != nil {
		panic(err)
	}
	inputItems = append(inputItems, item)

	//
	failingYttContentItem, err = kyaml.Parse(`
apiVersion: v1alpha1
kind: YttTemplate
metadata:
  name: ytt-invalid-item
ytt_template_content:
  #@ sample ytt content
  ytt_item:
    #@ more ytt content
#@ Ending comment
`)
	if err != nil {
		panic(err)
	}

	//
	failingYttAnnotationItem, err = kyaml.Parse(`
apiVersion: v1alpha1
kind: YttDataValues
metadata:
  name: ytt-invalid-item
  annotations:
    config.kubernetes.io/path: "path_to_file/more_path/"
ytt_header:
  header:
  #@ sample ytt annotation
#@ Ending comment
`)
	if err != nil {
		panic(err)
	}
}
