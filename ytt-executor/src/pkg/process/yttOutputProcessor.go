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
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

// UnmarshalYttOutput Parses ytt output into kyaml.RNode and writes it to provided items list under
// config.YttOutputElementKey
//
// Parameters:
//   - yttOutput: bytes.Buffer containing ytt binary output
//   - items: list of output RNodes to write ytt output to
//
// Returns:
//   - error: from parsing ytt output OR insufficient output RNodes available
func UnmarshalYttOutput(yttOutput bytes.Buffer, items []*kyaml.RNode) error {
	// Debug raw ytt output
	logger.LogDetailedDebug("Processing ytt binary output", map[string]string{
		"rawOutput": yttOutput.String(),
	})

	// Split bytes.buffer on "---"
	splitBuffer := bytes.Split(yttOutput.Bytes(), []byte("---"))

	if len(items) <= 0 {
		return fmt.Errorf("no output file with kind: %s provided", config.YttOutputFileKind)
	}

	// Check counts of files / ytt output provided / available
	if len(splitBuffer) > len(items) {
		// Generate error if not enough output items made available
		logger.LogDetailedError("Ytt output required more files than available", map[string]string{
			"ytt_output_count":    strconv.Itoa(len(splitBuffer)),
			"provided_file_count": strconv.Itoa(len(items)),
		})

		return errors.New("ytt output contained more files than available")

		// Generate warning if too many output items made available
	} else if len(splitBuffer) < len(items) {
		logger.LogDetailedWarning("Ytt output had more files provided than needed", map[string]string{
			"ytt_output_count":    strconv.Itoa(len(splitBuffer)),
			"provided_file_count": strconv.Itoa(len(items)),
		})
	}

	// Unmarshal and assemble each item into []kyaml.RNode
	for i, yttBytePart := range splitBuffer {
		if items[i].Field(config.YttOutputElementKey) == nil {
			logger.LogError(fmt.Sprintf(
				"Output file: %s, did not contain required output key: %s",
				items[i].GetAnnotations()["config.kubernetes.io/path"],
				config.YttOutputElementKey,
			))
			return fmt.Errorf(
				"output file: %s, did not contain required output key: %s",
				items[i].GetAnnotations()["config.kubernetes.io/path"],
				config.YttOutputElementKey,
			)
		}

		// Parse byteBuffer to kyaml.RNode
		dataPart, err := kyaml.Parse(string(yttBytePart))
		if err != nil {
			logger.LogDetailedError("Failed to parse ytt output item", map[string]string{
				"item_index":  strconv.Itoa(i),
				"output_dump": string(yttBytePart),
				"full_output": yttOutput.String(),
			})
			return err
		}

		// Generate info message depending on action (write / overwrite)
		if items[i].Field(config.YttOutputElementKey).Value.IsNilOrEmpty() {
			logger.LogInfo(fmt.Sprintf(
				"Overwriting file: %s, %s key",
				items[i].GetAnnotations()["config.kubernetes.io/path"],
				config.YttOutputElementKey,
			))
		} else {
			logger.LogInfo(fmt.Sprintf(
				"Writing to file: %s, %s key",
				items[i].GetAnnotations()["config.kubernetes.io/path"],
				config.YttOutputElementKey,
			))
		}

		// Set field in output items
		err = items[i].PipeE(kyaml.SetField(config.YttOutputElementKey, dataPart))
		if err != nil {
			return err
		}
	}
	return nil
}
