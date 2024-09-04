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
	"strings"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

// LogLevel Default log level to use
var LogLevel = LogLevelInfo

// LogStack list of framework.Result items to output as kpt log
var LogStack framework.Results

// LogLevel enum as a string representation
var LogLevelStrings = []string{"DEBUG", "INFO", "WARNING", "ERROR"}

// logLevels enumerator for possible log levels to use
type logLevels int

const (
	LogLevelDebug logLevels = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// SetLogLevel changes log level to one of the following: LogLevelStrings
func SetLogLevel(logLevel string) {
	for i, levelString := range LogLevelStrings {
		if levelString == strings.ToUpper(logLevel) {
			LogLevel = logLevels(i)
		}
	}
}

// LogDetailed appends framework.Result to LogStack with given logLevel, message and nullable detailed map
//
// Parameters:
//   - level: filter log saving based on current LogLevel, level has to be greater or equal to current one.
//   - message: main message of the log
//   - detailed: extra map for details, nullable
func LogDetailed(level logLevels, message string, detailed map[string]string) {
	if level >= LogLevel {
		LogStack = append(LogStack, &framework.Result{
			Message:  message,
			Severity: framework.Severity(LogLevelStrings[level]),
			Tags:     detailed,
		})
	}
}

// Log call to LogDetailed with only level and message, leaving detailed = nil
//
// Parameters:
//   - level: filter log saving based on current LogLevel, level has to be greater or equal to current one.
//   - message: main message of the log
func Log(level logLevels, message string) {
	LogDetailed(level, message, nil)
}

// LogDebug shorter hand reference for debug Log
func LogDebug(message string) {
	Log(LogLevelDebug, message)
}

// LogDetailedDebug shorter hand reference for debug LogDetailed
func LogDetailedDebug(message string, detailed map[string]string) {
	LogDetailed(LogLevelDebug, message, detailed)
}

// LogInfo shorter hand reference for info Log
func LogInfo(message string) {
	Log(LogLevelInfo, message)
}

// LogDetailedInfo shorter hand reference for info LogDetailed
func LogDetailedInfo(message string, detailed map[string]string) {
	LogDetailed(LogLevelInfo, message, detailed)
}

// LogWarning shorter hand reference for warning Log
func LogWarning(message string) {
	Log(LogLevelWarning, message)
}

// LogDetailedWarning shorter hand reference for warning LogDetailed
func LogDetailedWarning(message string, detailed map[string]string) {
	LogDetailed(LogLevelWarning, message, detailed)
}

// LogError shorter hand reference for error Log
func LogError(message string) {
	Log(LogLevelError, message)
}

// LogDetailedError shorter hand reference for error LogDetailed
func LogDetailedError(message string, detailed map[string]string) {
	LogDetailed(LogLevelError, message, detailed)
}
