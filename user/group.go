// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type groupField int

// Field names for group database.
const (
	G_NAME groupField = 1 << iota
	G_PASSWD
	G_GID
	G_MEMBER

	G_ALL
)

// A Group represents the format of a group on the system.
type Group struct {
	// Group name. (Unique)
	Name string

	// Hashed password
	//
	// The (hashed) group password. If this field is empty, no password is needed.
	password string

	// The numeric group ID. (Unique)
	GID int

	// User list
	//
	// A list of the usernames that are members of this group, separated by commas.
	UserList []string

	// A system group?
	IsOfSystem bool
}

// AddGroup returns a new Group.
func NewGroup(name string, members ...string) *Group {
	return &Group{
		Name:     name,
		password: "",
		GID:      -1,
		UserList: members,
	}
}

// NewSystemGroup adds a system group.
func NewSystemGroup(name string, members ...string) *Group {
	return &Group{
		Name:     name,
		password: "",
		GID:      -1,
		UserList: members,

		IsOfSystem: true,
	}
}

func (g *Group) filename() string { return _GROUP_FILE }

func (g *Group) name() string { return g.Name }

func (g *Group) String() string {
	return fmt.Sprintf("%s:%s:%d:%s\n",
		g.Name, g.password, g.GID, strings.Join(g.UserList, ","))
}

// parseGroup parses the row of a group.
func parseGroup(row string) (*Group, error) {
	fields := strings.Split(row, ":")
	if len(fields) != 4 {
		return nil, ErrRow
	}

	gid, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, &fieldError{_GROUP_FILE, row, "GID"}
	}

	return &Group{
		Name:     fields[0],
		password: fields[1],
		GID:      gid,
		UserList: strings.Split(fields[3], ","),
	}, nil
}

// == Lookup
//

// lookUp parses the group line searching a value into the field.
// Returns nil if it is not found.
func (*Group) lookUp(line string, field, value interface{}) interface{} {
	_field := field.(groupField)
	allField := strings.Split(line, ":")
	arrayField := make(map[int][]string)
	intField := make(map[int]int)

	arrayField[3] = strings.Split(allField[3], ",")

	// Check integers
	var err error
	if intField[2], err = strconv.Atoi(allField[2]); err != nil {
		panic(&fieldError{_GROUP_FILE, line, "GID"})
	}

	// Check fields
	var isField bool
	if G_NAME&_field != 0 && allField[0] == value.(string) {
		isField = true
	} else if G_PASSWD&_field != 0 && allField[1] == value.(string) {
		isField = true
	} else if G_GID&_field != 0 && intField[2] == value.(int) {
		isField = true
	} else if G_MEMBER&_field != 0 && checkGroup(arrayField[3], value.(string)) {
		isField = true
	} else if G_ALL&_field != 0 {
		isField = true
	}

	if isField {
		return &Group{
			Name:     allField[0],
			password: allField[1],
			GID:      intField[2],
			UserList: arrayField[3],
		}
	}
	return nil
}

// LookupGID looks up a group by group ID.
func LookupGID(gid int) (*Group, error) {
	entries, err := LookupInGroup(G_GID, gid, 1)
	if err != nil {
		return nil, err
	}

	return entries[0], err
}

// LookupGroup looks up a group by name.
func LookupGroup(name string) (*Group, error) {
	entries, err := LookupInGroup(G_NAME, name, 1)
	if err != nil {
		return nil, err
	}

	return entries[0], err
}

// LookupInGroup looks up a group by the given values.
//
// The count determines the number of fields to return:
//   n > 0: at most n fields
//   n == 0: the result is nil (zero fields)
//   n < 0: all fields
func LookupInGroup(field groupField, value interface{}, n int) ([]*Group, error) {
	iEntries, err := lookUp(&Group{}, field, value, n)
	if err != nil {
		return nil, err
	}

	// == Convert to type group
	valueSlice := reflect.ValueOf(iEntries)
	entries := make([]*Group, valueSlice.Len())

	for i := 0; i < valueSlice.Len(); i++ {
		entries[i] = valueSlice.Index(i).Interface().(*Group)
	}

	return entries, err
}

// Getgroups returns a list of the numeric ids of groups that the caller
// belongs to.
func Getgroups() []int {
	user := GetUsername()
	list := make([]int, 0)

	// The user could have its own group.
	if g, err := LookupGroup(user); err == nil {
		list = append(list, g.GID)
	}

	groups, err := LookupInGroup(G_MEMBER, user, -1)
	if err != nil {
		if err != ErrNoFound {
			panic(err)
		}
	}

	for _, v := range groups {
		list = append(list, v.GID)
	}
	return list
}

// GetgroupsName returns a list of the groups that the caller belongs to.
func GetgroupsName() []string {
	user := GetUsername()
	list := make([]string, 0)

	// The user could have its own group.
	if _, err := LookupGroup(user); err == nil {
		list = append(list, user)
	}

	groups, err := LookupInGroup(G_MEMBER, user, -1)
	if err != nil {
		if err != ErrNoFound {
			panic(err)
		}
	}
	for _, v := range groups {
		list = append(list, v.Name)
	}

	return list
}

// == Add
//

// AddGroup adds a group.
func AddGroup(name string, members ...string) (gid int, err error) {
	g := NewGroup(name, members...)
	_, err = g.Add()
	if err != nil {
		return
	}

	// TODO: GShadow should have AddGShadow?
	gs := &GShadow{name, "", []string{""}, members}
	if err = gs.Add(nil); err != nil {
		return
	}

	return g.GID, nil
}

// AddSystemGroup adds a system group.
func AddSystemGroup(name string, members ...string) (gid int, err error) {
	g := NewSystemGroup(name, members...)
	_, err = g.Add()
	if err != nil {
		return
	}

	// TODO: like above
	gs := &GShadow{name, "", []string{""}, members}
	if err = gs.Add(nil); err != nil {
		return
	}

	return g.GID, nil
}

// Add adds a new group.
// Whether GID is < 0, it will choose the first id available in the range set
// in the system configuration.
func (g *Group) Add() (gid int, err error) {
	loadConfig()

	group, err := LookupGroup(g.Name)
	if err != nil && err != ErrNoFound {
		return 0, err
	}
	if group != nil {
		return 0, ErrExist
	}

	if g.Name == "" {
		return 0, RequiredError("Name")
	}

	var db *dbfile
	if g.GID < 0 {
		db, gid, err = nextGUID(g.IsOfSystem)
		if err != nil {
			db.close()
			return 0, err
		}
		g.GID = gid
	} else {
		db, err = openDBFile(_GROUP_FILE, os.O_WRONLY|os.O_APPEND)
		if err != nil {
			return
		}

		// Check if Id is unique.
		_, err = LookupGID(g.GID)
		if err == nil {
			return 0, IdUsedError(g.GID)
		} else if err != ErrNoFound {
			return 0, err
		}
	}

	g.password = "x"

	_, err = db.file.WriteString(g.String())
	err2 := db.close()
	if err2 != nil && err == nil {
		err = err2
	}
	return
}

// == Remove
//

// DelGroup removes a group from the system.
func DelGroup(name string) (err error) {
	err = del(name, &Group{})
	if err == nil {
		err = del(name, &GShadow{})
	}
	return
}

// == Change
//

// AddUsersToGroup adds the members to a group.
func AddUsersToGroup(name string, members ...string) error {
	if len(members) == 0 {
		return fmt.Errorf("no members to add")
	}
	for i, v := range members {
		if v == "" {
			return EmptyError(fmt.Sprintf("members[%s]", strconv.Itoa(i)))
		}
	}

	gr, err := LookupGroup(name)
	if err != nil {
		return err
	}

	// Check if some member is already in the file.
	for _, u := range gr.UserList {
		for _, m := range members {
			if u == m {
				return fmt.Errorf("user %q is already set", u)
			}
		}
	}

	if len(gr.UserList) == 1 && gr.UserList[0] == "" {
		gr.UserList = members
	} else {
		gr.UserList = append(gr.UserList, members...)
	}

	return edit(name, gr)
}

// == Utility
//

// checkGroup indicates if a value is into a group.
func checkGroup(group []string, value string) bool {
	for _, v := range group {
		if v == value {
			return true
		}
	}
	return false
}

// * * *

// EmptyError reports an empty field.
type EmptyError string

func (e EmptyError) Error() string { return "empty field: " + string(e) }
