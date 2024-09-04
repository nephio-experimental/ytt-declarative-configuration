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

package commandExec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
)

// ExecuteYttForTemplate executes ytt binary with provided arguments in current or provided WorkDirectory
//
// Parameters:
//   - yttArgs: array of arguments used by function should consist of {"-f", "FILE_NAME",...}
//
// Returns:
//   - bytes.Buffer: direct output of ytt binary execution
//   - error: from getting directory or failing to execute ytt binary
func ExecuteYttForTemplate(yttArgs []string) (output bytes.Buffer, err error) {
	command := exec.Command(config.YttBinaryName, yttArgs...)

	// Define variables
	var outputBuffer, errorBuffer bytes.Buffer
	var workDir string

	// Check config workDir for starting '/' or used current WorkDirectory
	if len(config.YttWorkDirectory) > 0 && config.YttWorkDirectory[0] == '/' {
		workDir = config.YttWorkDirectory

	} else {
		workDir, err = os.Getwd()
		if err != nil {
			return outputBuffer, err
		}
		workDir = config.YttWorkDirectory + workDir
	}

	// Set command directory, Stdout and Stderr
	command.Dir = workDir
	command.Stdout = &outputBuffer
	command.Stderr = &errorBuffer

	// Debug details
	logger.LogDetailedDebug("Executing ytt binary", map[string]string{
		"ytt_bin_name": config.YttBinaryName,
		"work_dir":     workDir,
		"args":         fmt.Sprintf("%+v", yttArgs),
	})

	// Execute command
	err = command.Run()

	// Throw error back for better feedback
	if err != nil {
		return errorBuffer, fmt.Errorf(
			"ytt: %s (stderr: %s)",
			err,
			errorBuffer.String(),
		)
	}

	// Return outputBuffer
	return outputBuffer, err
}

// CatFile reads files
func CatFile(fileName string) string {
	command := exec.Command("cat", fileName)

	// Define variables
	var outputBuffer bytes.Buffer
	var workDir string

	// Check config workDir for starting '/' or used current WorkDirectory
	if len(config.YttWorkDirectory) > 0 && config.YttWorkDirectory[0] == '/' {
		workDir = config.YttWorkDirectory
	} else {
		workDir, _ = os.Getwd()
		workDir = config.YttWorkDirectory + workDir
	}

	command.Dir = workDir
	command.Stdout = &outputBuffer

	_ = command.Run()
	return outputBuffer.String()
}
