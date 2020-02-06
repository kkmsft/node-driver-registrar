// +build windows

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
	"errors"
	"fmt"
	"os"
)

func Umask(mask int) (int, error) {
	return -1, errors.New("umask not supported in Windows")
}

func CleanupSocketFile(socketPath string) error {
	if _, err := os.Lstat(socketPath); err != nil {
		// If the file does not exist, then the cleanup can be considered successful.
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to lstat the socket %s with error: %+v", socketPath, err)
	}

	// TODO: Until the bug - https://github.com/golang/go/issues/33357 is fixed, os.stat wouldn't return the
	// right mode(socket) on windows. Hence deleting the file, without checking whether
	// its a socket, on windows.
	if err := os.Remove(socketPath); err != nil {
		return fmt.Errorf("failed to remove stale socket %s with error: %+v", socketPath, err)
	}
	return nil
}
