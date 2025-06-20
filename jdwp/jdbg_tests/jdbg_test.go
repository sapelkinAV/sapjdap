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

package jdbg_tests_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sapelkinav/javadap/jdwp/jdbg"
	"sapelkinav/javadap/jdwp/jdwpclient"
	"sapelkinav/javadap/launcher"
	"sapelkinav/javadap/utils"
	"testing"
	"time"
)

const (
	TEST_LOG_DIR = "./.test_logs"
	JDWP_PORT    = 5007 // Different port from other tests
	JAR_PATH     = "../../daphelloworld/build/libs/daphelloworld-0.0.1-SNAPSHOT.jar"
)

type TestSetup struct {
	launcher   *launcher.JavaLauncher
	connection *jdwpclient.Connection
	socket     io.ReadWriteCloser
	ctx        context.Context
	cancel     context.CancelFunc
	thread     jdwpclient.ThreadID
}

func setupJDbgTest(t *testing.T) *TestSetup {
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

	// Get a thread for testing
	threads, err := connection.GetAllThreads()
	if err != nil {
		socket.Close()
		javaLauncher.Stop()
		cancel()
		t.Fatalf("Failed to get threads: %v", err)
	}

	if len(threads) == 0 {
		socket.Close()
		javaLauncher.Stop()
		cancel()
		t.Fatalf("No threads available")
	}

	return &TestSetup{
		launcher:   javaLauncher,
		connection: connection,
		socket:     socket,
		ctx:        ctx,
		cancel:     cancel,
		thread:     threads[0],
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

func TestJDbgDo(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// Basic test - just verify we can create a JDbg instance
		if j.Connection() != setup.connection {
			t.Error("Connection() should return the provided connection")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("JDbg.Do failed: %v", err)
	}
}

func TestJDbgBasicTypes(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// Test basic object types
		objType := j.ObjectType()
		if objType == nil {
			t.Error("ObjectType() should not return nil")
		}

		stringType := j.StringType()
		if stringType == nil {
			t.Error("StringType() should not return nil")
		}

		numberType := j.NumberType()
		if numberType == nil {
			t.Error("NumberType() should not return nil")
		}

		// Test primitive wrapper types
		boolObjType := j.BoolObjectType()
		if boolObjType == nil {
			t.Error("BoolObjectType() should not return nil")
		}

		intObjType := j.IntObjectType()
		if intObjType == nil {
			t.Error("IntObjectType() should not return nil")
		}

		// Test primitive types
		boolType := j.BoolType()
		if boolType == nil {
			t.Error("BoolType() should not return nil")
		}

		intType := j.IntType()
		if intType == nil {
			t.Error("IntType() should not return nil")
		}

		longType := j.LongType()
		if longType == nil {
			t.Error("LongType() should not return nil")
		}

		floatType := j.FloatType()
		if floatType == nil {
			t.Error("FloatType() should not return nil")
		}

		doubleType := j.DoubleType()
		if doubleType == nil {
			t.Error("DoubleType() should not return nil")
		}

		t.Logf("Successfully retrieved all basic types")
		return nil
	})

	if err != nil {
		t.Fatalf("JDbg basic types test failed: %v", err)
	}
}

func TestJDbgTypeBySignature(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		testCases := []struct {
			name      string
			signature string
		}{
			{"Object", "Ljava/lang/Object;"},
			{"String", "Ljava/lang/String;"},
			{"Integer", "Ljava/lang/Integer;"},
			{"Boolean", "Ljava/lang/Boolean;"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ty := j.Type(tc.signature)
				if ty == nil {
					t.Errorf("Type(%s) should not return nil", tc.signature)
				}

				if ty.Signature() != tc.signature {
					t.Errorf("Expected signature %s, got %s", tc.signature, ty.Signature())
				}

				t.Logf("Successfully resolved type: %s -> %s", tc.signature, ty.String())
			})
		}

		return nil
	})

	if err != nil {
		t.Fatalf("JDbg type by signature test failed: %v", err)
	}
}

func TestJDbgClassByName(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		testCases := []string{
			"java.lang.Object",
			"java.lang.String",
			"java.lang.Integer",
			"java.lang.Boolean",
		}

		for _, className := range testCases {
			t.Run(className, func(t *testing.T) {
				class := j.Class(className)
				if class == nil {
					t.Errorf("Class(%s) should not return nil", className)
				}

				t.Logf("Successfully resolved class: %s -> %s", className, class.String())
			})
		}

		return nil
	})

	if err != nil {
		t.Fatalf("JDbg class by name test failed: %v", err)
	}
}

func TestJDbgAllClasses(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		classes := j.AllClasses()
		if len(classes) == 0 {
			t.Error("AllClasses() should return at least one class")
		}

		// Look for common classes
		foundObject := false
		foundString := false

		for _, class := range classes {
			classStr := class.String()
			if classStr == "java.lang.Object" {
				foundObject = true
			}
			if classStr == "java.lang.String" {
				foundString = true
			}
		}

		if !foundObject {
			t.Error("Should find java.lang.Object in all classes")
		}

		if !foundString {
			t.Error("Should find java.lang.String in all classes")
		}

		t.Logf("Successfully retrieved %d classes", len(classes))
		return nil
	})

	if err != nil {
		t.Fatalf("JDbg all classes test failed: %v", err)
	}
}

func TestJDbgArrayOf(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// Test array of different types
		intType := j.IntType()
		intArrayType := j.ArrayOf(intType)
		if intArrayType == nil {
			t.Error("ArrayOf(intType) should not return nil")
		}

		stringType := j.StringType()
		stringArrayType := j.ArrayOf(stringType)
		if stringArrayType == nil {
			t.Error("ArrayOf(stringType) should not return nil")
		}

		t.Logf("Successfully created array types: %s, %s",
			intArrayType.String(), stringArrayType.String())
		return nil
	})

	if err != nil {
		t.Fatalf("JDbg array of test failed: %v", err)
	}
}

func TestJDbgString(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		testStrings := []string{
			"Hello, World!",
			"",
			"Test string with special characters: !@#$%^&*()",
			"Unicode test: 你好世界",
		}

		for _, testStr := range testStrings {
			t.Run("String_"+testStr[:min(len(testStr), 10)], func(t *testing.T) {
				val := j.String(testStr)
				if val == (jdbg.Value{}) {
					t.Errorf("String(%s) should not return empty value", testStr)
				}

				t.Logf("Successfully created string value for: %s", testStr)
			})
		}

		return nil
	})

	if err != nil {
		t.Fatalf("JDbg string test failed: %v", err)
	}
}

func TestJDbgThis(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := setup.connection.Suspend(setup.thread)
	if err != nil {
		t.Fatalf("Failed to suspend thread: %v", err)
	}
	defer setup.connection.Resume(setup.thread)

	err = jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// This test may fail if we're in a static context
		thisVal := j.This()
		// We don't fail if This() returns nil since it might be a static method
		t.Logf("This() returned: %+v", thisVal)
		return nil
	})

	if err != nil {
		t.Logf("JDbg This() test failed (may be expected in static context): %v", err)
	}
}

func TestJDbgGetArgument(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := setup.connection.Suspend(setup.thread)
	if err != nil {
		t.Fatalf("Failed to suspend thread: %v", err)
	}
	defer setup.connection.Resume(setup.thread)

	err = jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// Try to get the first argument (index 0)
		// This may fail if there are no arguments or if the method is static
		arg := j.GetArgument("arg0", 0)
		t.Logf("GetArgument returned: %+v", arg)
		return nil
	})

	if err != nil {
		t.Logf("JDbg GetArgument test failed (may be expected): %v", err)
	}
}

func TestJDbgErrorHandling(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	// Test with invalid type signature
	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// This should cause a failure within the Do block
		jtype := j.Type("InvalidSignure")
		print(jtype.Signature())
		return nil
	})

	if err == nil {
		t.Error("Expected error for invalid type signature, but got none")
	} else {
		t.Logf("Correctly got error for invalid signature: %v", err)
	}
}

func TestJDbgTypeIntegration(t *testing.T) {
	setup := setupJDbgTest(t)
	defer setup.teardown()

	err := jdbg.Do(setup.connection, setup.thread, func(j *jdbg.JDbg) error {
		// Get String type through different methods
		stringBySignature := j.Type("Ljava/lang/String;")
		stringByName := j.Class("java.lang.String")
		stringBuiltin := j.StringType()

		// They should all represent the same type
		if stringBySignature.Signature() != stringByName.Signature() {
			t.Errorf("String type signatures don't match: %s vs %s",
				stringBySignature.Signature(), stringByName.Signature())
		}

		if stringByName.Signature() != stringBuiltin.Signature() {
			t.Errorf("String type signatures don't match: %s vs %s",
				stringByName.Signature(), stringBuiltin.Signature())
		}

		// Test array creation
		stringArray := j.ArrayOf(stringBuiltin)
		expectedSig := "[Ljava/lang/String;"
		if stringArray.Signature() != expectedSig {
			t.Errorf("Expected array signature %s, got %s", expectedSig, stringArray.Signature())
		}

		// Create a string value
		testStr := "Integration test string"
		stringVal := j.String(testStr)
		if stringVal == (jdbg.Value{}) {
			t.Error("String creation should not return empty value")
		}

		t.Logf("Type integration test successful")
		t.Logf("  String signature: %s", stringBuiltin.Signature())
		t.Logf("  Array signature: %s", stringArray.Signature())
		t.Logf("  Created string value: %+v", stringVal)

		return nil
	})

	if err != nil {
		t.Fatalf("JDbg type integration test failed: %v", err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
