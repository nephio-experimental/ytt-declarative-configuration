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
	"errors"
	"os"
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/commandExec"
	"github.com/stretchr/testify/assert"
)

// TestWriteToFileBasic a set of basic tests to confirm basic functionality of method
func TestWriteToFileBasic(t *testing.T) {
	// Setup TempDir for testing
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	// Function args
	type args struct {
		filePath string
		data     string
	}

	// Test structure
	tests := []struct {
		args            args             // Function arguments
		name            string           // Name of the test
		wantErr         bool             // Is error expected
		errorCheck      func(error) bool // Error check
		expectedContent string           // Written file confirmation
	}{ // Test list

		// Basic test, write to a new file
		{
			args{
				filePath: "test_file.txt",
				data:     "initial_data\n",
			},
			"Write to File",
			false,
			nil,
			"initial_data\n",
		},

		// Basic test write additional data to already existing file
		{
			args{
				filePath: "test_file.txt",
				data:     "amended_data",
			},
			"Write to existing File",
			false,
			nil,
			"initial_data\n---\namended_data",
		},

		// Basic test to write without file name to catch os.IsNotExist error
		{
			args{
				filePath: "",
				data:     "any_data",
			},
			"Write without file name",
			true,
			os.IsNotExist,
			"",
		},
	}
	for _, tt := range tests {
		// Loop through tests
		t.Run(tt.name, func(t *testing.T) {

			// Execute function
			err := WriteToFile(tt.args.filePath, tt.args.data)

			// Check if error exists
			if err != nil {
				// If error was not expected
				if !tt.wantErr {
					t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)

					// Check if error matched expectation
				} else if !tt.errorCheck(err) {
					t.Errorf("WriteToFile() error = %v, wantErr %v, errorCheck %v", err, tt.wantErr, tt.errorCheck(err))
				}

				// If error is nil but was expected
			} else if tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check if expected content was provided and then compare it
			if len(tt.expectedContent) > 0 {
				cat := commandExec.CatFile(tempDir + "/" + tt.args.filePath)
				assert.Equal(t, tt.expectedContent, cat)
			}
		})
	}
}

// Pseudo Mock function replacement for os.MkdirAll
func MkdirAll(string, os.FileMode) error {
	return errors.New("os.MkdirAll(string, os.FileMode) error")
}

// Pseudo Mock function replacement for (*bufio.Writer).WriteString
func WriteString(*bufio.Writer, string) (int, error) {
	return 0, errors.New("writer.WriteString(s) error")
}

// Pseudo Mock function replacement for (*bufio.Writer).Flush
func Flush(*bufio.Writer) error {
	return errors.New("writer.Flush() error")
}

// Pseudo Mock function replacement for (*os.File).Close
func FileClose(*os.File) error {
	return errors.New("file.Close() error")
}

// TestWriteToFileMock test various error points with mock like approach
func TestWriteToFileMock(t *testing.T) {
	// Setup TempDir for testing
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	// Test file.Close()
	t.Run("Test file.Close() error", func(t *testing.T) {
		fileClose = FileClose
		err := WriteToFile("FileClose.txt", "test_data")
		assert.Equal(t, errors.New("file.Close() error"), err)
	})

	// Test writer.Flush()
	t.Run("Test writer.Flush() error", func(t *testing.T) {
		writerFlush = Flush
		err := WriteToFile("flush.txt", "test_data")
		assert.Equal(t, errors.New("writer.Flush() error"), err)
	})

	// Test writer.WriteString(string)
	t.Run("Test writer.WriteString(string) error", func(t *testing.T) {
		writerWriteString = WriteString
		err := WriteToFile("writeString.txt", "test_data")
		assert.Equal(t, errors.New("writer.WriteString(s) error"), err)
	})

	// Test writer.WriteString(string) repeated for full coverage
	t.Run("Test writer.WriteString(string) error", func(t *testing.T) {
		writerWriteStringSeparator = WriteString
		err := WriteToFile("writeString.txt", "test_data")
		assert.Equal(t, errors.New("writer.WriteString(s) error"), err)
	})

	// Test os.MkdirAll(path, perm)
	t.Run("Test os.MkdirAll(path, perm) error", func(t *testing.T) {
		osMkdirAll = MkdirAll
		err := WriteToFile("somePath/mkdirAll.txt", "test_data")
		assert.Equal(t, errors.New("os.MkdirAll(string, os.FileMode) error"), err)
	})

}
