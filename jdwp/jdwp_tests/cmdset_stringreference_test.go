// Copyright (C) 2017 Google Inc.
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

package jdwp_tests_test

import (
	"fmt"
	"sapelkinav/javadap/jdwp/jdwpclient"
	"testing"
)

type StringID = jdwpclient.StringID

func TestGetString(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testStrings := []string{
		"Hello, World!",
		"",
		"Test string with special characters: !@#$%^&*()",
		"Unicode test: 你好世界",
		"Multi-line\nstring\ntest",
		"Tab\ttest",
		"Quote \"test\"",
		"Backslash \\ test",
		"Very long string: " + string(make([]byte, 1000)),
	}

	for _, testStr := range testStrings {
		t.Run("String_"+testStr[:min(len(testStr), 20)], func(t *testing.T) {
			stringID, err := setup.connection.CreateString(testStr)
			if err != nil {
				t.Fatalf("CreateString failed for '%s': %v", testStr, err)
			}

			if stringID == 0 {
				t.Errorf("Expected non-zero StringID for '%s'", testStr)
			}

			retrievedStr, err := setup.connection.GetString(stringID)
			if err != nil {
				t.Fatalf("GetString failed for StringID %d: %v", stringID, err)
			}

			if retrievedStr != testStr {
				t.Errorf("String mismatch: expected '%s', got '%s'", testStr, retrievedStr)
			}

			t.Logf("Successfully created and retrieved string: '%s' (ID: %d)", testStr, stringID)
		})
	}
}

func TestGetStringWithMultipleStrings(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	stringMap := make(map[string]int)
	testStrings := []string{
		"First string",
		"Second string",
		"Third string",
		"Fourth string",
		"Fifth string",
	}

	for _, testStr := range testStrings {
		stringID, err := setup.connection.CreateString(testStr)
		if err != nil {
			t.Fatalf("CreateString failed for '%s': %v", testStr, err)
		}
		stringMap[testStr] = int(stringID)
		t.Logf("Created string '%s' with ID %d", testStr, stringID)
	}

	for testStr, stringID := range stringMap {
		retrievedStr, err := setup.connection.GetString(StringID(stringID))
		if err != nil {
			t.Fatalf("GetString failed for StringID %d: %v", stringID, err)
		}

		if retrievedStr != testStr {
			t.Errorf("String mismatch for ID %d: expected '%s', got '%s'", stringID, testStr, retrievedStr)
		}

		t.Logf("Successfully retrieved string '%s' for ID %d", retrievedStr, stringID)
	}
}

func TestGetStringEdgeCases(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testCases := []struct {
		name string
		test string
	}{
		{"Empty string", ""},
		{"Single character", "A"},
		{"Null character", "\x00"},
		{"Only whitespace", "   "},
		{"Only newlines", "\n\n\n"},
		{"Mixed control chars", "\t\r\n\v\f"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stringID, err := setup.connection.CreateString(tc.test)
			if err != nil {
				t.Fatalf("CreateString failed for %s: %v", tc.name, err)
			}

			retrievedStr, err := setup.connection.GetString(stringID)
			if err != nil {
				t.Fatalf("GetString failed for %s (ID %d): %v", tc.name, stringID, err)
			}

			if retrievedStr != tc.test {
				t.Errorf("%s: expected '%q', got '%q'", tc.name, tc.test, retrievedStr)
			}

			t.Logf("%s: Successfully handled string '%q' (ID: %d)", tc.name, tc.test, stringID)
		})
	}
}

func TestGetStringUnicodeSupport(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	unicodeStrings := []string{
		"English: Hello",
		"Chinese: 你好",
		"Japanese: こんにちは",
		"Korean: 안녕하세요",
		"Arabic: مرحبا",
		"Russian: Привет",
		"Emoji: 🌍🚀✨",
		"Mixed: Hello 世界 🌎",
	}

	for _, testStr := range unicodeStrings {
		t.Run("Unicode_"+testStr[:min(len(testStr), 15)], func(t *testing.T) {
			stringID, err := setup.connection.CreateString(testStr)
			if err != nil {
				t.Fatalf("CreateString failed for unicode string '%s': %v", testStr, err)
			}

			retrievedStr, err := setup.connection.GetString(stringID)
			if err != nil {
				t.Fatalf("GetString failed for unicode StringID %d: %v", stringID, err)
			}

			if retrievedStr != testStr {
				t.Errorf("Unicode string mismatch: expected '%s', got '%s'", testStr, retrievedStr)
			}

			t.Logf("Successfully handled unicode string: '%s' (ID: %d)", testStr, stringID)
		})
	}
}

func TestStringReferenceIntegration(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	const numStrings = 10
	createdStrings := make(map[int]string)

	for i := 0; i < numStrings; i++ {
		testStr := fmt.Sprintf("Integration test string #%d", i)

		stringID, err := setup.connection.CreateString(testStr)
		if err != nil {
			t.Fatalf("CreateString failed for integration test string %d: %v", i, err)
		}

		createdStrings[int(stringID)] = testStr
		t.Logf("Created integration test string %d: '%s' (ID: %d)", i, testStr, stringID)
	}

	for stringID, expectedStr := range createdStrings {
		retrievedStr, err := setup.connection.GetString(StringID(stringID))
		if err != nil {
			t.Fatalf("GetString failed for integration test StringID %d: %v", stringID, err)
		}

		if retrievedStr != expectedStr {
			t.Errorf("Integration test failed for ID %d: expected '%s', got '%s'", stringID, expectedStr, retrievedStr)
		}
	}

	t.Logf("String reference integration test successful: created and retrieved %d strings", numStrings)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
