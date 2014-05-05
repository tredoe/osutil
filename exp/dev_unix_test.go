// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package exp

import (
	"testing"
)

func TestDev(t *testing.T) {
	devs, err := GetUSBremovables()
	if err != nil {
		t.Fatal(err)
	}

	_, err = FindPartition("foo", devs)
	if err != CmdFindPartError("foo") {
		t.Errorf("FindPartition should get an error")
	}
}
