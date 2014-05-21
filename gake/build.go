// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kless/goutil"
)

func Builder(pkg *Package) error {
	workDir, err := ioutil.TempDir("", "gake-")
	if err != nil {
		return err
	}
	goutil.AtExit(func() { os.RemoveAll(workDir) })

	// Copy all files to the temporary directory.
	for _, f := range pkg.files {
		src, err := ioutil.ReadFile(f.name)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(workDir, filepath.Base(f.name)), src, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
