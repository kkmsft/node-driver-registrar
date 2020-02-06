/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"k8s.io/klog"
)

func TestSocketPath(t *testing.T) {
	// Create in the path where the test is run
	socketPath := "reg.sock"
	fileExists := false

	// Test1: file does not exist. So clean up should be successful.
	err := CleanupSocketFile(socketPath)
	if err != nil {
		t.Fatalf("cleanup returned error: %+v", err)
	}

	// Test2: director does not exist. So clean up should be successful.
	err = CleanupSocketFile(filepath.Join("test", socketPath))
	if err != nil {
		t.Fatalf("cleanup returned error: %+v", err)
	}

	defer func() {
		if fileExists {
			os.Remove(socketPath)
		}
	}()

	// Test3: Create the file path and check the cleanup
	_, err = net.Listen("unix", socketPath)
	if err != nil {
		klog.Errorf("failed to listen on socket: %s with error: %+v", socketPath, err)
		os.Exit(1)
	}
	fileExists = true

	err = CleanupSocketFile(socketPath)
	if err != nil {
		t.Fatalf("cleanup returned error: %+v", err)
	}

	_, err = os.Lstat(socketPath)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("lstat error on file %s ", socketPath)
		}
	} else {
		t.Fatalf("socket file %s exists", socketPath)
	}
	fileExists = false

	// Test 4: Create a regular file
	f, err := os.Create(socketPath)
	if err != nil {
		t.Fatalf("create file failed: %s", socketPath)
	}
	f.Close()

	fileExists = true
	err = CleanupSocketFile(socketPath)
	if err != nil {
		t.Fatalf("cleanup returned error: %+v", err)
	}

	if runtime.GOOS == "windows" {
		_, err = os.Lstat(socketPath)
		if err != nil {
			if !os.IsNotExist(err) {
				t.Fatalf("lstat error on file %s ", socketPath)
			}
		} else {
			t.Fatalf("regular file %s exists", socketPath)
		}
		fileExists = false
	} else {
		_, err = os.Lstat(socketPath)
		if err != nil {
			if os.IsNotExist(err) {
				t.Fatalf("regular file %s got deleted", socketPath)
			} else {
				t.Fatalf("lstat error on file %s ", socketPath)
			}
		}
		fileExists = true
	}

}
