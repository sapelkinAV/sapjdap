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
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sapelkinav/javadap/jdwp/jdwpclient"
	"sapelkinav/javadap/launcher"
	"sapelkinav/javadap/utils"
	"testing"
	"time"
)

const (
	TEST_LOG_DIR = "./.test_logs"
	JDWP_PORT    = 5006
	JAR_PATH     = "../../daphelloworld/build/libs/daphelloworld-0.0.1-SNAPSHOT.jar"
)

type TestSetup struct {
	launcher   *launcher.JavaLauncher
	connection *jdwpclient.Connection
	socket     io.ReadWriteCloser
	ctx        context.Context
	cancel     context.CancelFunc
}

func setupJDWPTest(t *testing.T) *TestSetup {
	if err := utils.InitializeLogger(TEST_LOG_DIR, "debug"); err != nil {
		t.Fatalf("Failed to set up logger: %v", err)
	}

	jarPath, err := filepath.Abs(JAR_PATH)
	if err != nil {
		t.Fatalf("Failed to get absolute path for JAR: %v", err)
	}

	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		t.Skipf("JAR file not found at %s, skipping test", jarPath)
	}

	javaLauncher := launcher.NewJavaLauncher(jarPath, JDWP_PORT)
	if err := javaLauncher.Start(); err != nil {
		t.Fatalf("Failed to launch Java process: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	var socket io.ReadWriteCloser
	for i := 0; i < 10; i++ {
		if socket, err = net.Dial("tcp", fmt.Sprintf("localhost:%d", JDWP_PORT)); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if socket == nil {
		javaLauncher.Stop()
		cancel()
		t.Fatalf("Failed to connect to JDWP socket after 10 attempts: %v", err)
	}

	connection, err := jdwpclient.Open(ctx, socket)
	if err != nil {
		socket.Close()
		javaLauncher.Stop()
		cancel()
		t.Fatalf("Failed to open JDWP connection: %v", err)
	}

	return &TestSetup{
		launcher:   javaLauncher,
		connection: connection,
		socket:     socket,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (ts *TestSetup) teardown() {
	if ts.socket != nil {
		ts.socket.Close()
	}
	if ts.launcher != nil {
		ts.launcher.Stop()
	}
	if ts.cancel != nil {
		ts.cancel()
	}
}

func TestGetVersion(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	version, err := setup.connection.GetVersion()
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	if version.Description == "" {
		t.Error("Version description should not be empty")
	}

	if version.JDWPMajor <= 0 {
		t.Errorf("JDWP major version should be positive, got %d", version.JDWPMajor)
	}

	if version.JDWPMinor < 0 {
		t.Errorf("JDWP minor version should be non-negative, got %d", version.JDWPMinor)
	}

	if version.Version == "" {
		t.Error("JRE version should not be empty")
	}

	if version.Name == "" {
		t.Error("VM name should not be empty")
	}

	t.Logf("Version info - Description: %s, JDWP: %d.%d, JRE: %s, VM: %s",
		version.Description, version.JDWPMajor, version.JDWPMinor, version.Version, version.Name)
}

func TestGetIDSizes(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	idSizes, err := setup.connection.GetIDSizes()
	if err != nil {
		t.Fatalf("GetIDSizes failed: %v", err)
	}

	if idSizes.FieldIDSize <= 0 {
		t.Errorf("FieldIDSize should be positive, got %d", idSizes.FieldIDSize)
	}

	if idSizes.MethodIDSize <= 0 {
		t.Errorf("MethodIDSize should be positive, got %d", idSizes.MethodIDSize)
	}

	if idSizes.ObjectIDSize <= 0 {
		t.Errorf("ObjectIDSize should be positive, got %d", idSizes.ObjectIDSize)
	}

	if idSizes.ReferenceTypeIDSize <= 0 {
		t.Errorf("ReferenceTypeIDSize should be positive, got %d", idSizes.ReferenceTypeIDSize)
	}

	if idSizes.FrameIDSize <= 0 {
		t.Errorf("FrameIDSize should be positive, got %d", idSizes.FrameIDSize)
	}

	t.Logf("ID Sizes - Field: %d, Method: %d, Object: %d, ReferenceType: %d, Frame: %d",
		idSizes.FieldIDSize, idSizes.MethodIDSize, idSizes.ObjectIDSize,
		idSizes.ReferenceTypeIDSize, idSizes.FrameIDSize)
}

func TestGetAllClasses(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	classes, err := setup.connection.GetAllClasses()
	if err != nil {
		t.Fatalf("GetAllClasses failed: %v", err)
	}

	if len(classes) == 0 {
		t.Error("Expected at least one class to be loaded")
	}

	foundJavaLangObject := false
	for _, class := range classes {
		if class.Signature == "Ljava/lang/Object;" {
			foundJavaLangObject = true
			break
		}
	}

	if !foundJavaLangObject {
		t.Error("Expected to find java.lang.Object in loaded classes")
	}

	t.Logf("Found %d loaded classes", len(classes))
}

func TestGetClassesBySignature(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testCases := []struct {
		name      string
		signature string
		expectMin int
	}{
		{"Object class", "Ljava/lang/Object;", 1},
		{"String class", "Ljava/lang/String;", 1},
		{"Non-existent class", "Lcom/nonexistent/Class;", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			classes, err := setup.connection.GetClassesBySignature(tc.signature)
			if err != nil {
				t.Fatalf("GetClassesBySignature failed for %s: %v", tc.signature, err)
			}

			if len(classes) < tc.expectMin {
				t.Errorf("Expected at least %d classes for signature %s, got %d",
					tc.expectMin, tc.signature, len(classes))
			}

			for _, class := range classes {
				if class.Signature != tc.signature {
					t.Errorf("Expected signature %s, got %s", tc.signature, class.Signature)
				}
			}

			t.Logf("Found %d classes for signature %s", len(classes), tc.signature)
		})
	}
}

func TestGetAllThreads(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Error("Expected at least one thread to be running")
	}

	t.Logf("Found %d active threads", len(threads))
}

func TestCreateString(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	testStrings := []string{
		"Hello, World!",
		"",
		"Test string with special characters: !@#$%^&*()",
		"Unicode test: 你好世界",
	}

	for _, testStr := range testStrings {
		t.Run(fmt.Sprintf("String_%s", testStr), func(t *testing.T) {
			stringID, err := setup.connection.CreateString(testStr)
			if err != nil {
				t.Fatalf("CreateString failed for '%s': %v", testStr, err)
			}

			if stringID == 0 {
				t.Errorf("Expected non-zero StringID for '%s'", testStr)
			}

			t.Logf("Created string '%s' with ID %d", testStr, stringID)
		})
	}
}

func TestSuspendAndResumeAll(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	err := setup.connection.SuspendAll()
	if err != nil {
		t.Fatalf("SuspendAll failed: %v", err)
	}

	t.Log("Successfully suspended all threads")

	err = setup.connection.ResumeAll()
	if err != nil {
		t.Fatalf("ResumeAll failed: %v", err)
	}

	t.Log("Successfully resumed all threads")
}

func TestResumeAllExcept(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	if len(threads) == 0 {
		t.Skip("No threads available for ResumeAllExcept test")
	}

	err = setup.connection.SuspendAll()
	if err != nil {
		t.Fatalf("SuspendAll failed: %v", err)
	}

	testThread := threads[0]
	err = setup.connection.ResumeAllExcept(testThread)
	if err != nil {
		t.Fatalf("ResumeAllExcept failed: %v", err)
	}

	t.Logf("Successfully resumed all threads except thread %d", testThread)

	err = setup.connection.ResumeAll()
	if err != nil {
		t.Fatalf("Final ResumeAll failed: %v", err)
	}
}

func TestVMCommandsIntegration(t *testing.T) {
	setup := setupJDWPTest(t)
	defer setup.teardown()

	version, err := setup.connection.GetVersion()
	if err != nil {
		t.Fatalf("GetVersion failed: %v", err)
	}

	idSizes, err := setup.connection.GetIDSizes()
	if err != nil {
		t.Fatalf("GetIDSizes failed: %v", err)
	}

	classes, err := setup.connection.GetAllClasses()
	if err != nil {
		t.Fatalf("GetAllClasses failed: %v", err)
	}

	threads, err := setup.connection.GetAllThreads()
	if err != nil {
		t.Fatalf("GetAllThreads failed: %v", err)
	}

	stringID, err := setup.connection.CreateString("Integration test string")
	if err != nil {
		t.Fatalf("CreateString failed: %v", err)
	}

	t.Logf("Integration test successful - Version: %s, Classes: %d, Threads: %d, StringID: %d",
		version.Name, len(classes), len(threads), stringID)

	if len(classes) == 0 {
		t.Error("Expected at least one class")
	}

	if len(threads) == 0 {
		t.Error("Expected at least one thread")
	}

	if stringID == 0 {
		t.Error("Expected non-zero string ID")
	}

	if idSizes.ObjectIDSize <= 0 {
		t.Error("Expected positive object ID size")
	}
}
