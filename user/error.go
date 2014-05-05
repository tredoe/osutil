// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrExist   = errors.New("user or group already exists")
	ErrNoFound = errors.New("entry not found")
	ErrRow     = errors.New("format of row not valid")
	ErrSearch  = errors.New("no search")
)

// IsExist returns whether the error is known to report that an user or group
// already exists. It is satisfied by ErrExist.
func IsExist(err error) bool { return err == ErrExist }

// A fieldError records the file, line and field that caused the error at turning
// a field from string to int.
type fieldError struct {
	file  string
	line  string
	field string
}

func (e *fieldError) Error() string {
	return fmt.Sprintf("field %q on %s: could not be turned to int\n%s",
		e.field, e.file, e.line)
}

// An IdRangeError records an error during the search for a free id to use.
type IdRangeError struct {
	LastId   int
	IsSystem bool
	IsUser   bool
}

func (e *IdRangeError) Error() string {
	str := ""
	if e.IsSystem {
		str = "system "
	}
	if e.IsUser {
		str += "user: "
	} else {
		str += "group: "
	}
	str += strconv.Itoa(e.LastId)

	return "reached maximum identifier in " + str
}

// An HomeError reports an error at adding an account with invalid home directory.
type HomeError string

func (e HomeError) Error() string {
	return "invalid directory for the home directory of an account: " + string(e)
}

// An IdUsedError reports the presence of an identifier already used.
type IdUsedError int

func (e IdUsedError) Error() string { return "id used: " + strconv.Itoa(int(e)) }

// A RequiredError reports the name of a required field.
type RequiredError string

func (e RequiredError) Error() string { return "required field: " + string(e) }
