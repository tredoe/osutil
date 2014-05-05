// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build linux darwin

package exp

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kless/term"
)

// TODO: Check command "plymouth" in Mac; whether exists there and the arguments
// are the same than in Linux version.

const CMD_WRITE = "/bin/plymouth" // to write during graphical boot

var USE_CMD_WRITE bool

// ReadPassword reads a password directly from terminal or through a third program.
func ReadPassword(prompt string) (key []byte, err error) {
	if USE_CMD_WRITE {
		key, err = exec.Command(CMD_WRITE, "ask-for-password", "--prompt="+prompt).Output()
	} else {
		var n int
		pass := make([]byte, 40)

		n, err = term.ReadPassword(pass)
		if err == nil {
			key = pass[:n]
		}
	}

	if err != nil {
		return nil, fmt.Errorf("ReadPassword: %s", err)
	}
	return
}

// Writef prints a message using the program in CMD_WRITE or to Stderr.
func Writef(format string, a ...interface{}) {
	if USE_CMD_WRITE {
		exec.Command(CMD_WRITE, "message", "--text="+fmt.Sprintf(format, a...)).Run()
	} else {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// Writefln is like Writef, but adds a new line.
func Writefln(format string, a ...interface{}) { Writef(format+"\n", a) }
