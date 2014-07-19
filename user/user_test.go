// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"bufio"
	"io"
	"os"
	"testing"
)

func TestUserParser(t *testing.T) {
	f, err := os.Open(_USER_FILE)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			continue
		}

		if _, err = parseUser(string(line)); err != nil {
			t.Error(err)
		}
	}
}

func TestUserFull(t *testing.T) {
	entry, err := LookupUID(os.Getuid())
	if err != nil || entry == nil {
		t.Error(err)
	}

	entry, err = LookupUser("root")
	if err != nil || entry == nil {
		t.Error(err)
	}

	entries, err := LookupInUser(U_GID, 65534, -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInUser(U_GECOS, "", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInUser(U_DIR, "/bin", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInUser(U_SHELL, "/bin/false", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInUser(U_ALL, nil, -1)
	if err != nil || len(entries) == 0 {
		t.Error(err)
	}
}

func TestUserCount(t *testing.T) {
	count := 2
	entries, err := LookupInUser(U_SHELL, "/bin/false", count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}

	count = 5
	entries, err = LookupInUser(U_ALL, nil, count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}
}

func TestUserError(t *testing.T) {
	var err error

	if _, err = LookupUser("!!!???"); err != ErrNoFound {
		t.Error("expected to report ErrNoFound")
	}

	if _, err = LookupInUser(U_SHELL, "/bin/false", 0); err != ErrSearch {
		t.Error("expected to report ErrSearch")
	}

	u := &User{}
	if err = u.Add(); err != RequiredError("Name") {
		t.Error("expected to report RequiredError")
	}

	u = &User{Name: USER, Dir: config.useradd.HOME, Shell: config.useradd.SHELL}
	if err = u.Add(); err != HomeError(config.useradd.HOME) {
		t.Error("expected to report HomeError")
	}
}

func TestUser_Add(t *testing.T) {
	user := &User{Name: USER, UID: -1, GID: GID, Dir: "/tmp", Shell: "/bin/sh"}
	err := user.Add()
	if err != nil {
		t.Fatal(err)
	}
	if err = user.Add(); err == nil {
		t.Fatal("an user existent can not be added again")
	} else {
		if !IsExist(err) {
			t.Error("user: expected to report ErrExist")
		}
	}

	shadow := &Shadow{Name: USER, changed: 1, Max: 10, Warn: 10}
	if err = shadow.Add(nil); err != nil {
		t.Fatal(err)
	}
	if err = shadow.Add(nil); err == nil {
		t.Fatal("a shadowed user existent can not be added again")
	} else {
		if !IsExist(err) {
			t.Error("shadow: expected to report ErrExist")
		}
	}

	u, err := LookupUser(USER)
	if err != nil {
		t.Fatal(err)
	}
	s, err := LookupShadow(USER)
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != USER {
		t.Errorf("user: expected to get name %q", USER)
	}
	if s.Name != USER {
		t.Errorf("shadow: expected to get name %q", USER)
	}

	// System user

	user.Name = SYS_USER
	user.UID = -1
	user.GID = SYS_GID
	if err = user.AddSystem(); err != nil {
		t.Fatal(err)
	}

	shadow.Name = SYS_USER
	if err = shadow.Add(nil); err != nil {
		t.Fatal(err)
	}

	u, err = LookupUser(SYS_USER)
	if err != nil {
		t.Fatal(err)
	}
	s, err = LookupShadow(SYS_USER)
	if err != nil {
		t.Fatal(err)
	}

	if u.Name != SYS_USER {
		t.Errorf("system user: expected to get name %q", SYS_USER)
	}
	if s.Name != SYS_USER {
		t.Errorf("system shadow: expected to get name %q", SYS_USER)
	}
}
