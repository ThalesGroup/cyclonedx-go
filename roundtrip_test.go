// This file is part of CycloneDX Go
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0
// Copyright (c) OWASP Foundation. All Rights Reserved.

package cyclonedx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestRoundTripJSON(t *testing.T) {
	bomFilePaths, err := filepath.Glob("./testdata/*.json")
	require.NoError(t, err)

	for _, bomFilePath := range bomFilePaths {
		t.Run(filepath.Base(bomFilePath), func(t *testing.T) {
			// Read original BOM JSON
			inputFile, err := os.Open(bomFilePath)
			require.NoError(t, err)

			// Decode BOM
			var bom BOM
			require.NoError(t, NewBOMDecoder(inputFile, BOMFileFormatJSON).Decode(&bom))
			inputFile.Close()

			// Prepare encoding destination
			buf := bytes.Buffer{}

			// Encode BOM again
			err = NewBOMEncoder(&buf, BOMFileFormatJSON).
				SetPretty(true).
				Encode(&bom)
			require.NoError(t, err)

			// Sanity checks: BOM has to be valid
			assertValidBOM(t, buf.Bytes(), BOMFileFormatJSON, SpecVersion1_6)

			// Compare with snapshot
			assert.NoError(t, snapShooter.SnapshotMulti(filepath.Base(bomFilePath), buf.String()))
		})
	}
}

func TestRoundTripXML(t *testing.T) {
	bomFilePaths, err := filepath.Glob("./testdata/*.xml")
	require.NoError(t, err)

	for _, bomFilePath := range bomFilePaths {
		t.Run(filepath.Base(bomFilePath), func(t *testing.T) {
			// Read original BOM XML
			inputFile, err := os.Open(bomFilePath)
			require.NoError(t, err)

			// Decode BOM
			var bom BOM
			require.NoError(t, NewBOMDecoder(inputFile, BOMFileFormatXML).Decode(&bom))
			inputFile.Close()

			// Prepare encoding destination
			buf := bytes.Buffer{}

			// Encode BOM again
			err = NewBOMEncoder(&buf, BOMFileFormatXML).
				SetPretty(true).
				Encode(&bom)
			require.NoError(t, err)

			// Sanity check: BOM has to be valid
			assertValidBOM(t, buf.Bytes(), BOMFileFormatXML, SpecVersion1_6)

			// Compare with snapshot
			assert.NoError(t, snapShooter.SnapshotMulti(filepath.Base(bomFilePath), buf.String()))
		})
	}
}

// This test uses JSON sample files from official CycloneDX specification repo:
// https://github.com/CycloneDX/specification/tree/master/tools/src/test/resources
func TestJSONSchemasFromSpecRepo(t *testing.T) {
	root := "testdata"
	schemaMap := map[string]string{
		// JSON schema start with CycloneDX 1.2
		"1.2": "schema/bom-1.2.schema.json",
		"1.3": "schema/bom-1.3.schema.json",
		"1.4": "schema/bom-1.4.schema.json",
		"1.5": "schema/bom-1.5.schema.json",
		"1.6": "schema/bom-1.6.schema.json",
		// Add other mappings if needed
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Ignore XML files as this test focus on JSON
		if strings.HasSuffix(info.Name(), ".xml") {
			return nil
		}

		// Ignore .textproto protobuf files  as this test focus on JSON
		if strings.HasSuffix(info.Name(), ".textproto") {
			return nil
		}

		// Determine the version of the file
		version := filepath.Base(filepath.Dir(path))

		// Find the schema corresponding to the version
		schemaPath, ok := schemaMap[version]
		if !ok {
			return nil // Skip if no schema is found for the version
		}

		// Read original BOM JSON
		inputFile, err := os.Open(path)
		require.NoError(t, err)

		// Decode BOM
		var bom BOM
		require.NoError(t, NewBOMDecoder(inputFile, BOMFileFormatJSON).Decode(&bom))
		inputFile.Close()

		// Prepare encoding destination
		buf := bytes.Buffer{}

		// Encode BOM again
		err = NewBOMEncoder(&buf, BOMFileFormatJSON).
			SetPretty(true).
			Encode(&bom)
		require.NoError(t, err)

		if strings.HasPrefix(info.Name(), "valid-") && strings.HasSuffix(info.Name(), ".json") {
			t.Run(info.Name(), func(t *testing.T) {
				assert.True(t, validateJSONstream(&buf, schemaPath))
			})
		}

		if strings.HasPrefix(info.Name(), "invalid-") && strings.HasSuffix(info.Name(), ".json") {
			t.Run(info.Name(), func(t *testing.T) {
				assert.False(t, validateJSONstream(&buf, schemaPath))
			})
		}

		return nil
	})

	assert.NoError(t, err)
}

// For now this method is used in TestJSONSchemasFromSpecRepo
func validateJSONstream(reader io.Reader, schemaPath string) bool {
	// Load the schema
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	// Read the JSON data from the io.Reader
	jsonBytes, err := io.ReadAll(reader)
	if err != nil {
		return false
	}

	// Unmarshal the JSON data
	var jsonData interface{}
	if err = json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return false
	}

	// Create a document loader
	documentLoader := gojsonschema.NewGoLoader(jsonData)

	// Validate the JSON data against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false
	}

	return result.Valid()
}
