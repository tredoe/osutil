// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import "testing"

func TestDelUser(t *testing.T) {
	err := DelUser(USER)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = LookupUser(USER); err != ErrNoFound {
		t.Error("expected to get error: %s", ErrNoFound)
	}
	if _, err = LookupShadow(USER); err != ErrNoFound {
		t.Error("expected to get error: %s", ErrNoFound)
	}
}

func TestDelGroup(t *testing.T) {
	err := DelGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = LookupGroup(GROUP); err != ErrNoFound {
		t.Error("expected to get error: %s", ErrNoFound)
	}
	if _, err = LookupGShadow(GROUP); err != ErrNoFound {
		t.Error("expected to get error: %s", ErrNoFound)
	}
}

func TestZ(*testing.T) {
	removeTempFiles()
}
