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

package fileWriter

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
)

// Exposing functions for testing purposes
var (
	osMkdirAll                 = os.MkdirAll
	fileClose                  = (*os.File).Close
	writerWriteString          = (*bufio.Writer).WriteString
	writerWriteStringSeparator = (*bufio.Writer).WriteString
	writerFlush                = (*bufio.Writer).Flush
)

// WriteToFile writes given string data to given file
// If the file already exists data is appended with "---" separator
//
// Parameters:
//   - filePath: path for a given file that gets default work dir prepended to it
//   - data: string data to write to the file
//
// Returns:
//   - error: Any error that could be experienced when writing the file
func WriteToFile(filePath string, data string) error {
	fullPath := config.YttWorkDirectory + filePath

	// Check if given file already exists
	contentExists := true
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		contentExists = false
	}

	// Crate parent dir if file doesn't exist
	if !contentExists {
		if err := osMkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}
	}

	// open/create file for editing
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	// Create buffer writer
	writer := bufio.NewWriter(file)

	// Write divider if file existed
	if contentExists {
		if _, err := writerWriteStringSeparator(writer, "---\n"); err != nil {
			return err
		}
	}

	// Write data
	if _, err := writerWriteString(writer, data); err != nil {
		return err
	}

	// Flush writer
	if err := writerFlush(writer); err != nil {
		return err
	}

	// Exit successfully
	return fileClose(file)
}
