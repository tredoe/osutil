// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestCommand(t *testing.T) {
	CMD_MAIN := "test-gake"

	// Build the command
	err := exec.Command("go", "build", "-o", CMD_MAIN).Run()
	if err != nil {
		t.Fatal(err)
	}

	err = exec.Command("./"+CMD_MAIN, "testdata").Run()
	if err != nil {
		t.Fatal(err)
	}

	// Remove the commands
	for _, v := range []string{CMD_MAIN} {
		if err = os.Remove(v); err != nil {
			t.Log(err)
		}
	}
}
