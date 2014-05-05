// Copyright 2010 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package user

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/kless/osutil/file"
)

// A dbfile represents the database file.
type dbfile struct {
	sync.Mutex
	file *os.File
	rd   *bufio.Reader
}

// openDBFile opens a file.
func openDBFile(filename string, flag int) (*dbfile, error) {
	f, err := os.OpenFile(filename, flag, 0)
	if err != nil {
		return nil, err
	}

	db := &dbfile{file: f, rd: bufio.NewReader(f)}
	db.Lock()
	return db, nil
}

// close closes the file.
func (db *dbfile) close() error {
	db.Unlock()
	return db.file.Close()
}

// A structurer represents the structure of a row into a file.
type structurer interface {
	// lookUp is the parser to looking for a value in the field of given line.
	lookUp(line string, field, value interface{}) interface{}

	// filename returns the file name belongs to the file structure.
	filename() string

	String() string
}

// lookUp is a generic parser to looking for a value.
//
// The count determines the number of fields to return:
//   n > 0: at most n fields
//   n == 0: the result is nil (zero fields)
//   n < 0: all fields
func lookUp(structer structurer, field, value interface{}, n int) (interface{}, error) {
	if n == 0 {
		return nil, ErrSearch
	}

	dbf, err := openDBFile(structer.filename(), os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer dbf.close()

	// Lines where a field is matched.
	entries := make([]interface{}, 0, 0)

	for {
		line, _, err := dbf.rd.ReadLine()
		if err == io.EOF {
			break
		}

		entry := structer.lookUp(string(line), field, value)
		if entry != nil {
			entries = append(entries, entry)
		}

		if n < 0 {
			continue
		} else if n == len(entries) {
			break
		}
	}

	if len(entries) != 0 {
		return entries, nil
	}
	return nil, ErrNoFound
}

// edit is a generic editor for the given user/group name.
// If remove is true, it removes the structure of the user/group name.
//
// It is created a backup before of modify the original shadowed files.
//
// TODO: get better performance if start to store since when the file is edited.
// So there is to store the size of all lines read until that point to seek from
// there.
func edit(name string, struc structurer, remove bool) (err error) {
	// Backup
	filename := struc.filename()
	if filename == _SHADOW_FILE || filename == _GSHADOW_FILE {
		if err = file.Copy(filename, filename+"-"); err != nil {
			return err
		}
	}

	dbf, err := openDBFile(filename, os.O_RDWR)
	if err != nil {
		return err
	}
	defer func() {
		e := dbf.close()
		if e != nil && err == nil {
			err = e
		}
	}()

	var buf bytes.Buffer
	_name := []byte(name)
	isFound := false

	for {
		line, err2 := dbf.rd.ReadBytes('\n')
		if err2 == io.EOF {
			break
		}

		if !isFound && bytes.HasPrefix(line, _name) {
			isFound = true
			if remove { // skip user
				continue
			}

			line = []byte(struc.String())
		}
		if _, err = buf.Write(line); err != nil {
			return err
		}
	}

	if isFound {
		if _, err = dbf.file.Seek(0, os.SEEK_SET); err != nil {
			return
		}

		var n int
		n, err = dbf.file.Write(buf.Bytes())
		if err != nil {
			return
		}
		err = dbf.file.Truncate(int64(n))
	}

	return
}
