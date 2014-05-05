// Copyright 2012 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import "testing"

var (
	USER_KEY1 = []byte("123")
	USER_KEY2 = []byte("456")

	GROUP_KEY1 = []byte("abc")
	GROUP_KEY2 = []byte("def")
)

func TestUserCrypt(t *testing.T) {
	s, err := LookupShadow(USER)
	if err != nil {
		t.Fatal(err)
	}
	s.Passwd(USER_KEY1)
	if err = config.crypter.Verify(s.password, USER_KEY1); err != nil {
		t.Fatalf("expected to get the same hashed password for %q", USER_KEY1)
	}

	if err = ChPasswd(USER, USER_KEY2); err != nil {
		t.Fatalf("expected to change password: %s", err)
	}
	s, _ = LookupShadow(USER)
	if err = config.crypter.Verify(s.password, USER_KEY2); err != nil {
		t.Fatalf("ChPasswd: expected to get the same hashed password for %q", USER_KEY2)
	}
}

func TestGroupCrypt(t *testing.T) {
	gs, err := LookupGShadow(GROUP)
	if err != nil {
		t.Fatal(err)
	}
	gs.Passwd(GROUP_KEY1)
	if err = config.crypter.Verify(gs.password, GROUP_KEY1); err != nil {
		t.Fatalf("expected to get the same hashed password for %q", GROUP_KEY1)
	}

	if err = ChGPasswd(GROUP, GROUP_KEY2); err != nil {
		t.Fatalf("expected to change password: %s", err)
	}
	gs, _ = LookupGShadow(GROUP)
	if err = config.crypter.Verify(gs.password, GROUP_KEY2); err != nil {
		t.Fatalf("ChGPasswd: expected to get the same hashed password for %q", GROUP_KEY2)
	}
}

func TestUserLock(t *testing.T) {
	err := LockUser(USER)
	if err != nil {
		t.Fatal(err)
	}
	s, err := LookupShadow(USER)
	if err != nil {
		t.Fatal(err)
	}
	if s.password[0] != _LOCK_CHAR {
		t.Fatalf("expected to get password starting with '%c', got: '%c'",
			_LOCK_CHAR, s.password[0])
	}

	err = UnlockUser(USER)
	if err != nil {
		t.Fatal(err)
	}
	s, err = LookupShadow(USER)
	if err != nil {
		t.Fatal(err)
	}
	if s.password[0] == _LOCK_CHAR {
		t.Fatalf("no expected to get password starting with '%c'", _LOCK_CHAR)
	}
}

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
