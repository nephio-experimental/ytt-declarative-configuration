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

package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var logStackIndex int = 0

func TestSetLogLevel(t *testing.T) {
	// Function arguments
	type args struct {
		logLevel string
	}

	// Test structure
	tests := []struct {
		name        string
		expectedLog logLevels
		args        args
	}{ // Test List

		// Set log level to 'Error'
		{
			"Set log level to Error",
			LogLevelError,
			args{"Error"},
		},

		// Set log level to 'WARNING'
		{
			"Set log level to WARNING",
			LogLevelWarning,
			args{"WARNING"},
		},

		// Set log level to 'info'
		{
			"Set log level to info",
			LogLevelInfo,
			args{"info"},
		},

		// Set log level to 'dEbUg'
		{
			"Set log level to dEbUg",
			LogLevelDebug,
			args{"dEbUg"},
		},

		// Set log Level to 'security'
		{
			"Fail to change log level (security)",
			LogLevelDebug,
			args{"security"},
		},
	}

	// Loop through tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute function
			SetLogLevel(tt.args.logLevel)

			// Check if log Level was set successfully
			assert.Equal(t, tt.expectedLog, LogLevel)
		})
	}
}

func TestLog(t *testing.T) {
	// Function arguments
	type args struct {
		message string
	}

	// Test structure
	tests := []struct {
		name          string
		logType       func(string)
		expectedLevel logLevels
		args          args
	}{ // Test List

		{
			"Log error",
			LogError,
			LogLevelError,
			args{"Basic ERROR log"},
		},

		{
			"Log warning",
			LogWarning,
			LogLevelWarning,
			args{"Basic WARNING log"},
		},

		{
			"Log info",
			LogInfo,
			LogLevelInfo,
			args{"Basic INFO log"},
		},

		{
			"Log debug",
			LogDebug,
			LogLevelDebug,
			args{"Basic DEBUG log"},
		},
	}

	// Test loop
	for _, tt := range tests {
		if t.Run(tt.name, func(t *testing.T) {
			tt.logType(tt.args.message)

			// Compare expected logLevel to LogLevel saved
			assert.Equal(t, LogLevelStrings[tt.expectedLevel], string(LogStack[logStackIndex].Severity))

			// Compare message expected to message in LogStack
			assert.Equal(t, tt.args.message, LogStack[logStackIndex].Message)

			// Check total logs collected to be +1 to current loop
			assert.Equal(t, logStackIndex+1, len(LogStack))
		}) {
			logStackIndex++
		}
	}
}

func TestLogDetailed(t *testing.T) {
	// Function arguments
	type args struct {
		message  string
		detailed map[string]string
	}

	// Test structure
	tests := []struct {
		name          string
		logType       func(string, map[string]string)
		expectedLevel logLevels
		args          args
	}{ // Test List

		{
			"Log detailed error",
			LogDetailedError,
			LogLevelError,
			args{
				"Detailed ERROR log",
				map[string]string{"Additional_info": "Error description"},
			},
		},

		{
			"Log detailed warning",
			LogDetailedWarning,
			LogLevelWarning,
			args{
				"Detailed WARNING log",
				map[string]string{"Additional_info": "Warning description"},
			},
		},

		{
			"Log detailed info",
			LogDetailedInfo,
			LogLevelInfo,
			args{
				"Detailed INFO log",
				map[string]string{"Additional_info": "Info description"},
			},
		},

		{
			"Log detailed Debug",
			LogDetailedDebug,
			LogLevelDebug,
			args{
				"Detailed DEBUG log",
				map[string]string{"Additional_info": "Debug description"},
			},
		},
	}
	for _, tt := range tests {
		if t.Run(tt.name, func(t *testing.T) {
			tt.logType(tt.args.message, tt.args.detailed)

			// Compare expected logLevel to LogLevel saved
			assert.Equal(t, LogLevelStrings[tt.expectedLevel], string(LogStack[logStackIndex].Severity))

			// Compare message expected to message in LogStack
			assert.Equal(t, tt.args.message, LogStack[logStackIndex].Message)

			// Check details
			assert.Equal(t, tt.args.detailed, LogStack[logStackIndex].Tags)

			// Check total logs collected to be +1 to current loop
			assert.Equal(t, logStackIndex+1, len(LogStack))
		}) {
			logStackIndex++
		}
	}
}

// TODO: add a test to capture logLevel filtering
