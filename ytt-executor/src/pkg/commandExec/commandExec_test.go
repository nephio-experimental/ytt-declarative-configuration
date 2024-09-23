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
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestCatFile(t *testing.T) {

	// Coverage test for CatFile for custom root based directory
	t.Run("Test workDir='/' and nonexistent_file.txt", func(t *testing.T) {
		config.YttWorkDirectory = "/"
		output := CatFile("nonexistent_file.txt")
		assert.Equal(t, "", output)
	})

	// Coverage test for CatFile for get os.pwd() coverage
	t.Run("Test with nonexistent_file.txt", func(t *testing.T) {
		config.YttWorkDirectory = ""
		output := CatFile("nonexistent_file.txt")
		assert.Equal(t, "", output)
	})
}

func TestExecuteYttForTemplate(t *testing.T) {
	// Test structure
	tests := []struct {
		name        string
		args        []string
		wantOutput  *bytes.Buffer
		wantErr     bool
		preTestFunc func()
	}{ // Test list

		// Execute successfully with root based workDirectory and echo hello to work with any env
		{
			"Test ExecuteYttForTemplate with workDir='/' and echo hello",
			[]string{"hello"},
			bytes.NewBufferString("hello\n"),
			false,
			func() {
				config.YttBinaryName = "echo"
				config.YttWorkDirectory = "/"
			},
		},

		// Execute successfully echo hello again to work with any env
		{
			"Test ExecuteYttForTemplate with echo hello again",
			[]string{"hello", "again"},
			bytes.NewBufferString("hello again\n"),
			false,
			func() {
				config.YttWorkDirectory = ""
				config.YttBinaryName = "echo"
			},
		},

		// Execute ls on non_existent_dir to simulate ExecuteYttForTemplate error failure
		{
			"Test ExecuteYttForTemplate to fail using ls non_existent_dir",
			[]string{"non_existent_dir"},
			bytes.NewBufferString("ls: cannot access 'non_existent_dir': No such file or directory\n"),
			true,
			func() {
				config.YttBinaryName = "ls"
			},
		},

		//
		{
			"Test ExecuteYttForTemplate to fail using invalid binary",
			[]string{},
			&bytes.Buffer{},
			true,
			func() {
				config.YttBinaryName = "ytt2"
			},
		},
	}

	// Loop through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if tt.ConfigChange() function should be ran
			if tt.preTestFunc != nil {
				tt.preTestFunc()
			}

			// Execute function
			output, err := ExecuteYttForTemplate(tt.args)

			// Check if error received and not expected
			if err != nil && !tt.wantErr {
				t.Fatalf("error not expected %v", err)

				// Check if error expected but not received
			} else if err == nil && tt.wantErr {
				t.Fatalf("error expected but not received %v", err)
			}

			// Assert response
			assert.Equal(t, tt.wantOutput, &output)
		})
	}
}
