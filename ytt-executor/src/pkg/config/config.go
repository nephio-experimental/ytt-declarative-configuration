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

// Package config to store and modify easily accessible values for kpt-ytt template processing
package config

import (
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

// Default variables used in defining ytt strip-down functionality
// To be overridden by values provided in fnConfig using Configure
var (
	YttWorkDirectory           = ""                     // Directory to write ytt input files to (probably not needed)
	YttBinaryName              = "ytt"                  // Ytt binary name
	YttInputValuesFileHandling = ValuesIdentifierKind   // YttValuesIdentifier Enumerator to identify data-value-file handling
	YttInputValueFileKind      = "YttDataValues"        // Kind value to identify data-value-file
	YttNodeAnnotations         = "ytt_header"           // Yaml key to identify ytt annotation element
	YttNodeContent             = "ytt_template_content" // Yaml key to identify ytt content
	YttOutputFileHandling      = OutputFileKind         // YttOutputFileIdentifier Enumerator to identify ytt-output file handling
	YttOutputFileKind          = "Configuration"        //
	YttOutputElementKey        = "data"                 // Element key under which YTT output should be under
)

// Variables used by Package config to identify fnConfig fields to read when
// overriding default variables
var (
	configInputRootKey          = "input"          // Root node for input configuration
	configInputYttAnnotationKey = "ytt_header"     // Key used to identify ytt annotation element
	configInputYttContentKey    = "ytt_content"    // Key used to identify ytt content element
	configInputYttCiqIdentifier = "ciq_identifier" // Key used to identify ytt data-values-file identifier
	configOutputRootKey         = "output"         // Root node for output configuration
	configOutputKindKey         = "kind"           // Key used to identify output kind
	configOutputElementKey      = "output_key"     // Key used to identify output element key
	configDebugRootKey          = "debug"          // Root node for handling debug parameters
	configDebugWorkDirOverride  = "work_dir"       // Key used for overriding work directory
	configDebugYttBinOverride   = "bin_name"       // Key used for overriding binary name
	configDebugLogLevel         = "log_level"      // Key used for changing log level
)

// YttValuesIdentifier enumerator for identifying value-files handling
//
// ValuesIdentifierNone: do not identify value-files manually
//
// ValuesIdentifierKind: identify value-files based on kind
//
// ValuesIdentifierNamed: use specified files by name as value-files (WIP)
type YttValuesIdentifier int

const (
	ValuesIdentifierNone YttValuesIdentifier = iota
	ValuesIdentifierKind
)

// YttOutputFileIdentifier enumerator to identify output file handling
//
// OutputFileKind: determine output file by kind
type YttOutputFileIdentifier int

const (
	OutputFileKind YttOutputFileIdentifier = iota
)

// Configure parses fnConfig and overwrite default values for easily accessible go values
//
// Parameters:
//   - fnConfig: kyaml.RNode representing function config to be parsed, resourceList.FunctionConfig
func Configure(fnConfig *kyaml.RNode) error {

	// Input customization
	if inputs := fnConfig.Field(configInputRootKey); !inputs.IsNilOrEmpty() {
		inputs := inputs.Value

		// Check for user defined annotations key
		if !inputs.Field(configInputYttAnnotationKey).IsNilOrEmpty() {
			value, err := inputs.GetString(configInputYttAnnotationKey)
			if err != nil {
				return err
			}
			YttNodeAnnotations = value
		}

		// Check for user defined content key
		if !inputs.Field(configInputYttContentKey).IsNilOrEmpty() {
			value, err := inputs.GetString(configInputYttContentKey)
			if err != nil {
				return err
			}
			YttNodeContent = value
		}

		// Check for user defined ciq identifier
		if ciqIdentifier := inputs.Field(configInputYttCiqIdentifier); !ciqIdentifier.IsNilOrEmpty() {
			ciqIdentifier := ciqIdentifier.Value
			if !ciqIdentifier.Field("kind").IsNilOrEmpty() {
				YttInputValuesFileHandling = ValuesIdentifierKind
				value, err := ciqIdentifier.GetString("kind")
				if err != nil {
					return err
				}
				YttInputValueFileKind = value
			}
		}
	}

	// Output customization
	if outputs := fnConfig.Field(configOutputRootKey); !outputs.IsNilOrEmpty() {
		outputs := outputs.Value

		// Check for kind identifier
		if !outputs.Field(configOutputKindKey).IsNilOrEmpty() {
			value, err := outputs.GetString(configOutputKindKey)
			if err != nil {
				return err
			}
			YttOutputFileKind = value
		}

		// Check for output element key
		if !outputs.Field(configOutputElementKey).IsNilOrEmpty() {
			value, err := outputs.GetString(configOutputElementKey)
			if err != nil {
				return err
			}
			YttOutputElementKey = value
		}
	}

	// Debug customization
	if debug := fnConfig.Field(configDebugRootKey); !debug.IsNilOrEmpty() {
		debug := debug.Value

		// Check for work dir override (to facilitate non-container usage)
		if !debug.Field(configDebugWorkDirOverride).IsNilOrEmpty() {
			value, err := debug.GetString(configDebugWorkDirOverride)
			if err != nil {
				return err
			}
			YttWorkDirectory = value
		}

		// Check for ytt bin name override (to facilitate non-container usage)
		if !debug.Field(configDebugYttBinOverride).IsNilOrEmpty() {
			value, err := debug.GetString(configDebugYttBinOverride)
			if err != nil {
				return err
			}
			YttBinaryName = value
		}

		// Change debug level to given value
		if !debug.Field(configDebugLogLevel).IsNilOrEmpty() {
			debugLevel, err := debug.GetString(configDebugLogLevel)
			if err != nil {
				return err
			}
			logger.SetLogLevel(debugLevel)
		}
	}
	return nil
}
