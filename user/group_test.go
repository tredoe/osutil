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

func TestGroupParser(t *testing.T) {
	f, err := os.Open(_GROUP_FILE)
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

		if _, err = parseGroup(string(line)); err != nil {
			t.Error(err)
		}
	}
}

func TestGroupFull(t *testing.T) {
	entry, err := LookupGID(os.Getgid())
	if err != nil || entry == nil {
		t.Error(err)
	}

	entry, err = LookupGroup("root")
	if err != nil || entry == nil {
		t.Error(err)
	}

	entries, err := LookupInGroup(G_MEMBER, "", -1)
	if err != nil || entries == nil {
		t.Error(err)
	}

	entries, err = LookupInGroup(G_ALL, nil, -1)
	if err != nil || len(entries) == 0 {
		t.Error(err)
	}
}

func TestGroupCount(t *testing.T) {
	count := 5
	entries, err := LookupInGroup(G_ALL, nil, count)
	if err != nil || len(entries) != count {
		t.Error(err)
	}
}

func TestGroupError(t *testing.T) {
	var err error

	if _, err = LookupGroup("!!!???"); err != ErrNoFound {
		t.Error("expected to report ErrNoFound")
	}

	if _, err = LookupInGroup(G_MEMBER, "", 0); err != ErrSearch {
		t.Error("expected to report ErrSearch")
	}

	g := &Group{}
	if err = g.Add(); err != RequiredError("Name") {
		t.Error("expected to report RequiredError")
	}
}

func TestGetGroups(t *testing.T) {
	gids := Getgroups()
	gnames := GetgroupsName()

	for i, gid := range gids {
		g, err := LookupGID(gid)
		if err != nil {
			t.Error(err)
		}

		if g.Name != gnames[i] {
			t.Errorf("expected to match GID and group name")
		}
	}
}

func TestGroup_Add(t *testing.T) {
	member0 := "m0"
	member1 := "m1"

	var err error
	GID, err = AddGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = AddGroup(GROUP); err == nil {
		t.Fatal("a group existent can not be added again")
	} else {
		if !IsExist(err) {
			t.Error("expected to report ErrExist")
		}
	}

	g, err := LookupGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}
	sg, err := LookupGShadow(GROUP)
	if err != nil {
		t.Fatal(err)
	}

	if g.Name != GROUP {
		t.Errorf("group: expected to get name %q", GROUP)
	}
	if sg.Name != GROUP {
		t.Errorf("sgroup: expected to get name %q", GROUP)
	}

	// System group

	if SYS_GID, err = AddSystemGroup(SYS_GROUP, member0, member1); err != nil {
		t.Fatal(err)
	}

	g, err = LookupGroup(SYS_GROUP)
	if err != nil {
		t.Fatal(err)
	}
	sg, err = LookupGShadow(SYS_GROUP)
	if err != nil {
		t.Fatal(err)
	}

	if g.Name != SYS_GROUP {
		t.Errorf("system group: expected to get name %q", SYS_GROUP)
	}
	if sg.Name != SYS_GROUP {
		t.Errorf("system sgroup: expected to get name %q", SYS_GROUP)
	}

	if g.UserList[0] != member0 || g.UserList[1] != member1 {
		t.Error("system group: expected to get members: %s", g.UserList)
	}
	if sg.UserList[0] != member0 || sg.UserList[1] != member1 {
		t.Error("system group: expected to get members: %s", sg.UserList)
	}
}

func TestGroup_Change(t *testing.T) {
	g_first, err := LookupGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}

	err = AddUsersToGroup(GROUP, USER, SYS_USER)
	if err != nil {
		t.Fatal(err)
	}

	g_last, err := LookupGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}

	if len(g_first.UserList) == len(g_last.UserList) ||
		g_last.UserList[0] != USER ||
		g_last.UserList[1] != SYS_USER {
		t.Error("expected to add users into a group")
	}
}
