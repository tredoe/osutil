// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Build uses the tool "go build" to compile the make files.
// Returns the working directory and the error, if any.
func Build(pkg *Package) (workDir string, err error) {
	workDir, err = ioutil.TempDir("", "gake-")
	if err != nil {
		return
	}

	// Copy all files to the temporary directory.
	for _, f := range pkg.files {
		src, err := ioutil.ReadFile(f.name)
		if err != nil {
			return "", err
		}

		err = ioutil.WriteFile(filepath.Join(workDir, filepath.Base(f.name)), src, 0644)
		if err != nil {
			return "", err
		}
	}

	if err = os.Chdir(workDir); err != nil {
		return "", err
	}

	cmd := exec.Command("go", "build")
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}

	return
}
