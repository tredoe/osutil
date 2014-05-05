// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build linux

package exp

import (
	"os"
	"os/exec"
)

func init() {
	// Check if there is an the external command to write.
	_, err := os.Stat(CMD_WRITE)
	if os.IsNotExist(err) {
		return
	}

	err = exec.Command(CMD_WRITE, "--ping").Run()
	if _, ok := err.(*exec.ExitError); !ok {

		//	if _, ok, _ := shutil.Run(CMD_WRITE + " --ping"); ok {
		USE_CMD_WRITE = true
	}
}
