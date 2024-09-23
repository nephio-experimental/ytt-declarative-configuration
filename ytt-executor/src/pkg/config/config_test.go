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

package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	"github.com/stretchr/testify/assert"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestConfigure(t *testing.T) {
	// Parse raw fnConfig.yaml and use that
	fnConfig, err := kyaml.Parse(`
input:
  ciq_identifier:
    kind: CustomYttDataValues
  ytt_header: custom_ytt_header
  ytt_content: custom_ytt_template_content
output:
  kind: CustomCNSConfigurationFiles
  output_key: custom_data
debug:
  work_dir: subDir
  bin_name: echo
  log_level: Debug
`)

	// Catch error from parsing config
	if err != nil {
		t.Fatalf("failed to parse sample fnConfig: %v", err)
	}

	// Read config and change values to al the required fields
	t.Run("Read config and assert values", func(t *testing.T) {
		// Execute configure
		err := Configure(fnConfig)
		if err != nil {
			t.Fatalf("Encountered error while reading fnConfig: %v", err)
		}

		// Assert fields
		assert.Equal(t, "CustomYttDataValues", YttInputValueFileKind)
		assert.Equal(t, "custom_ytt_header", YttNodeAnnotations)
		assert.Equal(t, "custom_ytt_template_content", YttNodeContent)
		assert.Equal(t, "CustomCNSConfigurationFiles", YttOutputFileKind)
		assert.Equal(t, "custom_data", YttOutputElementKey)
		assert.Equal(t, "subDir", YttWorkDirectory)
		assert.Equal(t, "echo", YttBinaryName)
		assert.Equal(t, logger.LogLevel, logger.LogLevelDebug)
	})
}

func TestConfigureErrors(t *testing.T) {
	// Test structure
	tests := []struct {
		name           string
		errorFieldName string
		errorFieldType string
		fnConfig       string
	}{ // Test list

		// Test catch error of parse inputs.ciq_identifier.kind
		{
			"Test fail to parse inputs.ciq_identifier.kind",
			"kind",
			"map[child_element:to_break_parsing]",
			`
input:
  ciq_identifier:
    kind:
      child_element: to_break_parsing
`,
		},

		// Test catch error of parse inputs.ytt_header
		{
			"Test fail to parse inputs.ytt_header",
			"ytt_header",
			"map[child_element:to_break_parsing]",
			`
input:
  ytt_header:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse inputs.ytt_content
		{
			"Test fail to parse inputs.ytt_content",
			"ytt_content",
			"map[child_element:to_break_parsing]",
			`
input:
  ytt_content:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse output.kind
		{
			"Test fail to parse output.kind",
			"kind",
			"map[child_element:to_break_parsing]",
			`
output:
  kind:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse output.output_key
		{
			"Test fail to parse output.output_key",
			"output_key",
			"map[child_element:to_break_parsing]",
			`
output:
  output_key:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse debug.work_dir
		{
			"Test fail to parse debug.work_dir",
			"work_dir",
			"map[child_element:to_break_parsing]",
			`
debug:
  work_dir:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse debug.bin_name
		{
			"Test fail to parse debug.bin_name",
			"bin_name",
			"map[child_element:to_break_parsing]",
			`
debug:
  bin_name:
    child_element: to_break_parsing
`,
		},

		// Test catch error of parse debug.log_level
		{
			"Test fail to parse debug.log_level",
			"log_level",
			"map[child_element:to_break_parsing]",
			`
debug:
  log_level:
    child_element: to_break_parsing
`,
		},
	}

	// Loop through tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse string into kyaml.RNode object
			fnConfig, err := kyaml.Parse(tt.fnConfig)
			if err != nil {
				t.Fatalf("malformed test input, unable to parse: %v", err)
			}

			// Execute function
			err = Configure(fnConfig)
			if err == nil {
				t.Fatalf("error was expected, name: %v, fnConfigSlice: %v", tt.name, tt.fnConfig)
			}

			// Compare error expected vs received
			assert.Equal(
				t,
				errors.New(fmt.Sprintf(
					"node %v is not a string: %v",
					tt.errorFieldName,
					tt.errorFieldType,
				)),
				err,
			)
		})
	}
}
