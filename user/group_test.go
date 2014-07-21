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
	if _, err = g.Add(); err != RequiredError("Name") {
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
	group := NewGroup(GROUP, MEMBERS...)
	_testGroup_Add(t, group, MEMBERS, false)

	group = NewSystemGroup(SYS_GROUP, MEMBERS...)
	_testGroup_Add(t, group, MEMBERS, true)
}

func _testGroup_Add(t *testing.T, group *Group, members []string, ofSystem bool) {
	prefix := "group"
	if ofSystem {
		prefix = "system " + prefix
	}

	id, err := group.Add()
	if err != nil {
		t.Fatal(err)
	}
	if id == -1 {
		t.Errorf("%s: got UID = -1", prefix)
	}

	if _, err = group.Add(); err == nil {
		t.Fatalf("%s: an existent group can not be added again", prefix)
	} else {
		if !IsExist(err) {
			t.Errorf("%s: expected to report ErrExist", prefix)
		}
	}

	if ofSystem {
		if !group.IsOfSystem() {
			t.Errorf("%s: IsOfSystem(): expected true")
		}
	} else {
		if group.IsOfSystem() {
			t.Errorf("%s: IsOfSystem(): expected false")
		}
	}

	// Check value stored

	name := ""
	if ofSystem {
		name = SYS_GROUP
	} else {
		name = GROUP
	}

	g, err := LookupGroup(name)
	if err != nil {
		t.Fatalf("%s: ", err)
	}

	if g.Name != name {
		t.Errorf("%s: expected to get name %q", prefix, name)
	}
	if g.UserList[0] != members[0] || g.UserList[1] != members[1] {
		t.Error("%s: expected to get members: %s", prefix, g.UserList)
	}
}

func TestGroup_Change(t *testing.T) {
	g_first, err := LookupGroup(GROUP)
	if err != nil {
		t.Fatal(err)
	}

	err = AddUsersToGroup(GROUP, "m0")
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
