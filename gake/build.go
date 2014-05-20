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

func Builder(files []string) error {
	workDir, err := ioutil.TempDir("", "gake")
	if err != nil {
		return err
	}
	goutil.AtExit(func() { os.RemoveAll(workDir) })

	// Copy all files to the temporary directory.
	for _, f := range files {
		src, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		if err = ioutil.WriteFile(filepath.Join(workDir, f), src, 0644); err != nil {
			return err
		}
	}

	return nil
}
