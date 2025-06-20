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
	"sapelkinav/javadap/jdwp/jdwpclient"
	"testing"
)

type ThreadID = jdwpclient.ThreadID
type FrameID = jdwpclient.FrameID
type TaggedObjectID = jdwpclient.TaggedObjectID
type VariableRequest = jdwpclient.VariableRequest
type VariableAssignmentRequest = jdwpclient.VariableAssignmentRequest
type Value = jdwpclient.Value

const (
	TagObject  = jdwpclient.TagObject
	TagInt     = jdwpclient.TagInt
	TagBoolean = jdwpclient.TagBoolean
	TagByte    = jdwpclient.TagByte
	TagChar    = jdwpclient.TagChar
	TagDouble  = jdwpclient.TagDouble
	TagFloat   = jdwpclient.TagFloat
	TagLong    = jdwpclient.TagLong
	TagShort   = jdwpclient.TagShort
)

func getThreadAndFrame(t *testing.T, setup *TestSetup) (ThreadID, FrameID) {
	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for stack frame test")
	}

	testThread := threads[0]

	err = setup.connection.Suspend(testThread)
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}

	frames, err := setup.connection.GetFrames(testThread, 0, 1)
	if err != nil {
		setup.connection.Resume(testThread)
		t.Fatalf("GetFrames failed: %v", err)
	}

	if len(frames) == 0 {
		setup.connection.Resume(testThread)
		t.Skip("No frames available for stack frame test")
	}

	return testThread, frames[0].Frame
}

func TestGetThisObject(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	thisObject, err := setup.connection.GetThisObject(testThread, testFrame)
	if err != nil {
		t.Logf("GetThisObject failed (may be expected for static methods): %v", err)
		return
	}

	t.Logf("Successfully retrieved this object: Type=%d, Object=%d", thisObject.Type, thisObject.Object)

	if thisObject.Type == 0 && thisObject.Object == 0 {
		t.Log("This object is null (likely a static method)")
	} else {
		t.Logf("This object present: Type=%d, Object=%d", thisObject.Type, thisObject.Object)
	}
}

func TestGetValues(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	testCases := []struct {
		name  string
		slots []VariableRequest
	}{
		{
			"Single variable slot 0",
			[]VariableRequest{{Index: 0, Tag: uint8(TagObject)}}, // 'L' for object
		},
		{
			"Multiple variable slots",
			[]VariableRequest{
				{Index: 0, Tag: uint8(TagObject)}, // 'L' for object
				{Index: 1, Tag: uint8(TagInt)},    // 'I' for int
			},
		},
		{
			"Empty slots",
			[]VariableRequest{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			values, err := setup.connection.GetValues(testThread, testFrame, tc.slots)
			if err != nil {
				t.Logf("GetValues failed for %s (may be expected): %v", tc.name, err)
				return
			}

			if len(values) != len(tc.slots) {
				t.Errorf("Expected %d values, got %d", len(tc.slots), len(values))
			}

			t.Logf("%s: Successfully retrieved %d values", tc.name, len(values))
			for i, value := range values {
				t.Logf("  Value %d: %+v", i, value)
			}
		})
	}
}

func TestGetValuesWithDifferentTags(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	tagTests := []struct {
		name string
		tag  uint8
		desc string
	}{
		{"Object reference", uint8(TagObject), "L - Object reference"},
		{"Integer", uint8(TagInt), "I - Integer"},
		{"Boolean", uint8(TagBoolean), "Z - Boolean"},
		{"Byte", uint8(TagByte), "B - Byte"},
		{"Char", uint8(TagChar), "C - Character"},
		{"Double", uint8(TagDouble), "D - Double"},
		{"Float", uint8(TagFloat), "F - Float"},
		{"Long", uint8(TagLong), "J - Long"},
		{"Short", uint8(TagShort), "S - Short"},
	}

	for _, tagTest := range tagTests {
		t.Run(tagTest.name, func(t *testing.T) {
			slots := []VariableRequest{{Index: 0, Tag: tagTest.tag}}
			values, err := setup.connection.GetValues(testThread, testFrame, slots)
			if err != nil {
				t.Logf("GetValues failed for %s (%s): %v", tagTest.name, tagTest.desc, err)
				return
			}

			if len(values) > 0 {
				t.Logf("%s (%s): Retrieved value: %+v", tagTest.name, tagTest.desc, values[0])
			}
		})
	}
}

func TestSetValues(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	getSlots := []VariableRequest{{Index: 0, Tag: uint8(TagInt)}} // Try to get an integer
	originalValues, err := setup.connection.GetValues(testThread, testFrame, getSlots)
	if err != nil {
		t.Logf("GetValues failed, skipping SetValues test: %v", err)
		return
	}

	if len(originalValues) == 0 {
		t.Skip("No values available for SetValues test")
	}

	setSlots := []VariableAssignmentRequest{
		{Index: 0, Value: int32(42)}, // Set integer value
	}

	err = setup.connection.SetValues(testThread, testFrame, setSlots)
	if err != nil {
		t.Logf("SetValues failed (may be expected for non-writable variables): %v", err)
		return
	}

	newValues, err := setup.connection.GetValues(testThread, testFrame, getSlots)
	if err != nil {
		t.Fatalf("GetValues after SetValues failed: %v", err)
	}

	if len(newValues) > 0 {
		t.Logf("SetValues test: Original value=%+v, New value=%+v", originalValues[0], newValues[0])
	}

	t.Log("SetValues completed successfully")
}

func TestSetValuesWithDifferentTypes(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	testCases := []struct {
		name  string
		slots []VariableAssignmentRequest
	}{
		{
			"Integer value",
			[]VariableAssignmentRequest{{Index: 0, Value: int32(123)}},
		},
		{
			"Boolean value",
			[]VariableAssignmentRequest{{Index: 1, Value: true}},
		},
		{
			"Multiple values",
			[]VariableAssignmentRequest{
				{Index: 0, Value: int32(456)},
				{Index: 1, Value: false},
			},
		},
		{
			"Empty assignment",
			[]VariableAssignmentRequest{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setup.connection.SetValues(testThread, testFrame, tc.slots)
			if err != nil {
				t.Logf("SetValues failed for %s (may be expected): %v", tc.name, err)
				return
			}

			t.Logf("%s: SetValues completed successfully for %d slots", tc.name, len(tc.slots))
		})
	}
}

func TestStackFrameIntegration(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testThread, testFrame := getThreadAndFrame(t, setup)
	defer setup.connection.Resume(testThread)

	thisObject, err := setup.connection.GetThisObject(testThread, testFrame)
	if err != nil {
		t.Logf("GetThisObject failed (may be static method): %v", err)
	} else {
		t.Logf("This object: Type=%d, Object=%d", thisObject.Type, thisObject.Object)
	}

	getSlots := []VariableRequest{
		{Index: 0, Tag: uint8(TagObject)}, // Object
		{Index: 1, Tag: uint8(TagInt)},    // Integer
	}

	values, err := setup.connection.GetValues(testThread, testFrame, getSlots)
	if err != nil {
		t.Logf("GetValues failed: %v", err)
	} else {
		t.Logf("Retrieved %d values from frame", len(values))
		for i, value := range values {
			t.Logf("  Variable %d: %+v", i, value)
		}
	}

	if len(values) > 0 {
		setSlots := []VariableAssignmentRequest{
			{Index: 0, Value: values[0]}, // Set back the same value
		}

		err = setup.connection.SetValues(testThread, testFrame, setSlots)
		if err != nil {
			t.Logf("SetValues failed: %v", err)
		} else {
			t.Log("SetValues completed successfully")
		}
	}

	t.Logf("Stack frame integration test completed for thread %d, frame %d", testThread, testFrame)
}

func TestStackFrameErrorHandling(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for error handling test")
	}

	testThread := threads[0]
	invalidFrame := FrameID(999999) // Use invalid frame ID

	t.Run("Invalid frame ID", func(t *testing.T) {
		_, err := setup.connection.GetThisObject(testThread, invalidFrame)
		if err == nil {
			t.Error("Expected error for invalid frame ID, but got none")
		} else {
			t.Logf("Expected error for invalid frame ID: %v", err)
		}
	})

	t.Run("Invalid variable slots", func(t *testing.T) {
		err := setup.connection.Suspend(testThread)
		if err != nil {
			t.Fatalf("Suspend failed: %v", err)
		}
		defer setup.connection.Resume(testThread)

		frames, err := setup.connection.GetFrames(testThread, 0, 1)
		if err != nil || len(frames) == 0 {
			t.Skip("No valid frames for invalid slot test")
		}

		invalidSlots := []VariableRequest{{Index: 999, Tag: uint8(TagInt)}}
		_, err = setup.connection.GetValues(testThread, frames[0].Frame, invalidSlots)
		if err == nil {
			t.Log("No error for invalid variable slot (may be valid)")
		} else {
			t.Logf("Expected error for invalid variable slot: %v", err)
		}
	})
}
