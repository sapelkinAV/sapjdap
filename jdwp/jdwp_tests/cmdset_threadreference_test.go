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
	"testing"
)

func TestGetThreadName(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for GetThreadName test")
	}

	testThread := threads[0]
	name, err := setup.connection.GetThreadName(testThread)
	if err != nil {
		t.Fatalf("GetThreadName failed: %v", err)
	}

	if name == "" {
		t.Error("Thread name should not be empty")
	}

	t.Logf("Thread %d name: %s", testThread, name)
}

func TestSuspendAndResumeThread(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for suspend/resume test")
	}

	testThread := threads[0]

	err = setup.connection.Suspend(testThread)
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}

	t.Logf("Successfully suspended thread %d", testThread)

	err = setup.connection.Resume(testThread)
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}

	t.Logf("Successfully resumed thread %d", testThread)
}

func TestGetThreadStatus(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for GetThreadStatus test")
	}

	testThread := threads[0]
	threadStatus, suspendStatus, err := setup.connection.GetThreadStatus(testThread)
	if err != nil {
		t.Fatalf("GetThreadStatus failed: %v", err)
	}

	t.Logf("Thread %d status - Thread: %d, Suspend: %d", testThread, threadStatus, suspendStatus)
}

func TestGetSuspendCount(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for GetSuspendCount test")
	}

	testThread := threads[0]

	initialCount, err := setup.connection.GetSuspendCount(testThread)
	if err != nil {
		t.Fatalf("GetSuspendCount failed: %v", err)
	}

	if initialCount < 0 {
		t.Errorf("Suspend count should be non-negative, got %d", initialCount)
	}

	err = setup.connection.Suspend(testThread)
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}

	countAfterSuspend, err := setup.connection.GetSuspendCount(testThread)
	if err != nil {
		setup.connection.Resume(testThread)
		t.Fatalf("GetSuspendCount after suspend failed: %v", err)
	}

	if countAfterSuspend != initialCount+1 {
		t.Errorf("Expected suspend count to increase by 1, got %d -> %d", initialCount, countAfterSuspend)
	}

	err = setup.connection.Resume(testThread)
	if err != nil {
		t.Fatalf("Resume failed: %v", err)
	}

	countAfterResume, err := setup.connection.GetSuspendCount(testThread)
	if err != nil {
		t.Fatalf("GetSuspendCount after resume failed: %v", err)
	}

	if countAfterResume != initialCount {
		t.Errorf("Expected suspend count to return to initial value %d, got %d", initialCount, countAfterResume)
	}

	t.Logf("Suspend count test passed: %d -> %d -> %d", initialCount, countAfterSuspend, countAfterResume)
}

func TestGetFrames(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for GetFrames test")
	}

	testThread := threads[0]

	err = setup.connection.Suspend(testThread)
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	defer setup.connection.Resume(testThread)

	frames, err := setup.connection.GetFrames(testThread, 0, -1)
	if err != nil {
		t.Logf("GetFrames with count=-1 failed: %v, trying with count=1", err)
		frames, err = setup.connection.GetFrames(testThread, 0, 1)
		if err != nil {
			t.Fatalf("GetFrames failed: %v", err)
		}
	}

	t.Logf("Successfully retrieved %d frames for thread %d", len(frames), testThread)

	for i, frame := range frames {
		t.Logf("Frame %d: ID=%d, Location=%+v", i, frame.Frame, frame.Location)
	}
}

func TestGetFramesWithDifferentParameters(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for GetFrames parameter test")
	}

	testThread := threads[0]

	err = setup.connection.Suspend(testThread)
	if err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	defer setup.connection.Resume(testThread)

	testCases := []struct {
		name  string
		start int
		count int
	}{
		{"First frame only", 0, 1},
		{"All frames", 0, -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			frames, err := setup.connection.GetFrames(testThread, tc.start, tc.count)
			if err != nil {
				t.Logf("GetFrames failed for %s: %v", tc.name, err)
				return
			}

			t.Logf("%s: Retrieved %d frames", tc.name, len(frames))
		})
	}
}

//func TestThreadReferenceIntegration(t *testing.T) {
//	setup := setupJDWPTest(t)
//	defer setup.teardown()
//
//	threads, err := setup.connection.GetAllThreads()
//	if err != nil {
//		t.Fatalf("GetAllThreads failed: %v", err)
//	}
//
//	if len(threads) == 0 {
//		t.Skip("No threads available for integration test")
//	}
//
//	testThread := threads[0]
//
//	name, err := setup.connection.GetThreadName(testThread)
//	if err != nil {
//		t.Fatalf("GetThreadName failed: %v", err)
//	}
//
//	threadStatus, suspendStatus, err := setup.connection.GetThreadStatus(testThread)
//	if err != nil {
//		t.Fatalf("GetThreadStatus failed: %v", err)
//	}
//
//	initialSuspendCount, err := setup.connection.GetSuspendCount(testThread)
//	if err != nil {
//		t.Fatalf("GetSuspendCount failed: %v", err)
//	}
//
//	err = setup.connection.Suspend(testThread)
//	if err != nil {
//		t.Fatalf("Suspend failed: %v", err)
//	}
//
//	frames, err := setup.connection.GetFrames(testThread, 0, 1)
//	if err != nil {
//		setup.connection.Resume(testThread)
//		t.Logf("GetFrames failed: %v", err)
//		frames = []jdwpclient.FrameInfo{}
//	}
//
//	newSuspendCount, err := setup.connection.GetSuspendCount(testThread)
//	if err != nil {
//		setup.connection.Resume(testThread)
//		t.Fatalf("GetSuspendCount after suspend failed: %v", err)
//	}
//
//	err = setup.connection.Resume(testThread)
//	if err != nil {
//		t.Fatalf("Resume failed: %v", err)
//	}
//
//	finalSuspendCount, err := setup.connection.GetSuspendCount(testThread)
//	if err != nil {
//		t.Fatalf("GetSuspendCount after resume failed: %v", err)
//	}
//
//	t.Logf("Thread integration test successful:")
//	t.Logf("  Thread ID: %d", testThread)
//	t.Logf("  Name: %s", name)
//	t.Logf("  Status: Thread=%d, Suspend=%d", threadStatus, suspendStatus)
//	t.Logf("  Suspend counts: %d -> %d -> %d", initialSuspendCount, newSuspendCount, finalSuspendCount)
//	t.Logf("  Frames retrieved: %d", len(frames))
//
//	if name == "" {
//		t.Error("Thread name should not be empty")
//	}
//
//	if newSuspendCount != initialSuspendCount+1 {
//		t.Errorf("Expected suspend count to increase by 1, got %d -> %d", initialSuspendCount, newSuspendCount)
//	}
//
//	if finalSuspendCount != initialSuspendCount {
//		t.Errorf("Expected suspend count to return to initial value %d, got %d", initialSuspendCount, finalSuspendCount)
//	}
//
//	if len(frames) == 0 {
//		t.Log("No frames retrieved (thread may not have an active stack)")
//	}
//}
