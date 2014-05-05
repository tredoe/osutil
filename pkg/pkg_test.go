// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkg

import "testing"

func TestPackager(t *testing.T) {
	sys := New(Deb)
	cmd := "curl"

	err := sys.Install(cmd)
	if err != nil {
		t.Errorf("\n%s", err)
	}

	if err = sys.Remove(false, cmd); err != nil {
		t.Errorf("\n%s", err)
	}
}
