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
	"testing"

	"github.com/nephio-experimental/ytt-declarative-configuration/pkg/logger"
	"github.com/stretchr/testify/assert"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestUnmarshalYttOutput(t *testing.T) {
	// Setup output list
	var outputList []*kyaml.RNode

	// Parse and add first item
	outputItem1, err := kyaml.Parse(`
apiVersion: v1alpha1
kind: OutputKind
metadata:
  name: output-1
data: "existing_data"
`)
	if err != nil {
		t.Fatalf("malformed test input, unable to parse: %v", err)
	}
	outputList = append(outputList, outputItem1)

	// Parse and add second item
	outputItem2, err := kyaml.Parse(`
apiVersion: v1alpha1
kind: OutputKind
metadata:
  name: output-2
other_data: some_info
data:
`)
	if err != nil {
		t.Fatalf("malformed test input, unable to parse: %v", err)
	}
	outputList = append(outputList, outputItem2)

	// Single yaml file output by ytt against outputList with 2 available items
	t.Run("Single bytes.Buffer output", func(t *testing.T) {
		// Copy output list
		outputCopy := make([]*kyaml.RNode, len(outputList))
		copy(outputCopy, outputList)

		// Create sample output bytes.Buffer
		sampleOutput := bytes.Buffer{}
		sampleOutput.WriteString("yttOutputKey: yttOutputElement\n")

		// Execute function
		err := UnmarshalYttOutput(sampleOutput, outputCopy)
		if err != nil {
			t.Fatalf("error not expected: %v", err)
		}

		// Check first item
		data, err := outputCopy[0].Pipe(kyaml.Get("data"))
		if err != nil {
			t.Fatalf("error not expected: %v", err)
		}
		assert.Equal(t, sampleOutput.String(), data.MustString())

		// Check second item to be empty / nil
		data, _ = outputCopy[1].Pipe(kyaml.Get("data"))
		if !data.IsNilOrEmpty() {
			t.Fatal("data item should not exist in secondary output.")
		}

		// Check for generated warning
		log := logger.LogStack[len(logger.LogStack)-2]
		assert.Equal(t, "Ytt output had more files provided than needed", log.Message)
	})

	// Happy test when ytt output matches file count
	t.Run("Matching bytes.Buffer output", func(t *testing.T) {
		// Copy output list
		outputCopy := make([]*kyaml.RNode, len(outputList))
		copy(outputCopy, outputList)

		// Create sample output bytes.Buffer
		sampleOutput := bytes.Buffer{}
		sampleOutput.WriteString(`
yttOutputKey1: yttOutputElement1
---
yttOutputKey2: yttOutputElement2
`)

		// Execute function
		err := UnmarshalYttOutput(sampleOutput, outputCopy)
		if err != nil {
			t.Fatalf("error not expected: %v", err)
		}

		// Check first item
		data, err := outputCopy[0].Pipe(kyaml.Get("data"))
		if err != nil {
			t.Fatalf("error not expected: %v", err)
		}
		assert.Equal(t, "yttOutputKey1: yttOutputElement1\n", data.MustString())

		// Check second item
		data, err = outputCopy[1].Pipe(kyaml.Get("data"))
		if err != nil {
			t.Fatalf("error not expected: %v", err)
		}
		assert.Equal(t, "yttOutputKey2: yttOutputElement2\n", data.MustString())

		// Check final 2 info messages to have correct feedback
		assert.Equal(
			t,
			"Overwriting file: , data key",
			logger.LogStack[len(logger.LogStack)-1].Message,
		)
		assert.Equal(
			t,
			"Writing to file: , data key",
			logger.LogStack[len(logger.LogStack)-2].Message,
		)
	})

	// Test when ytt output yields more output than files provided
	t.Run("Insufficient output for bytes.Buffer", func(t *testing.T) {
		// Copy output list
		outputCopy := make([]*kyaml.RNode, len(outputList))
		copy(outputCopy, outputList)

		// Create sample output bytes.Buffer
		sampleOutput := bytes.Buffer{}
		sampleOutput.WriteString(`
yttOutputKey1: yttOutputElement1
---
yttOutputKey2: yttOutputElement2
---
yttOutputKey3: yttOutputElement3
`)

		// Execute function
		err := UnmarshalYttOutput(sampleOutput, outputCopy)

		// Check error
		assert.Equal(
			t,
			errors.New("ytt output contained more files than available"),
			err,
		)
	})

	// Fail when trying to parse ytt output into yaml format (This is imporbable to occur with successful ytt output)
	t.Run("Invalid yaml yttOutput", func(t *testing.T) {
		// Copy output list
		outputCopy := make([]*kyaml.RNode, len(outputList))
		copy(outputCopy, outputList)

		// Create sample output bytes.Buffer
		sampleOutput := bytes.Buffer{}
		sampleOutput.WriteString("invalid: yaml: item")

		// Execute function
		err := UnmarshalYttOutput(sampleOutput, outputCopy)

		// Check error
		assert.Equal(
			t,
			errors.New("yaml: mapping values are not allowed in this context"),
			err,
		)
	})

	// Test fail when no key is present in the output file
	t.Run("Fail on no output element", func(t *testing.T) {
		// Copy output list
		testList := []*kyaml.RNode{kyaml.MustParse(`
apiVersion: v1alpha1
kind: OutputKind
metadata:
  name: output-without-key
some: unrelated_data
`)}

		// Create sample output bytes.Buffer
		sampleOutput := bytes.Buffer{}
		sampleOutput.WriteString("yttOutputKey: yttOutputElement\n")

		// Execute function
		err := UnmarshalYttOutput(sampleOutput, testList)

		// Check error
		assert.Equal(
			t,
			errors.New("output file: , did not contain required output key: data"),
			err,
		)
	})
}
