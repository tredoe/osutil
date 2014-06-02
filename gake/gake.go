// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"hash/adler32"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const (
	// SUBDIR_HOME is the directory where are stored the compiled programs
	SUBDIR_HOME = ".run"

	
	CMD_NAME = "bygake"
)

func main() {
	flag.Parse()

	// Get the home directory for the compiled programs
	HOME := os.Getenv(ENV_HOME)
	if HOME == "" {
		// In Unix systems, the environmentvariable is not set during boot init.
		if runtime.GOOS != "windows" {
			user, err := user.Current()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			} else {
				if user.Uid == "0" { // root
					HOME = "/root"
				}
			}
		}
		if HOME == "" {
			fmt.Fprintf(os.Stderr, "environment variable %s is not set\n", ENV_HOME)
			os.Exit(1)
		}
	}
	HOME = filepath.Join(HOME, SUBDIR_HOME)

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}

	errList := make([]error, 0)

	for i, arg := range args {
		absPath, err := filepath.Abs(arg)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		// Create the directory, if it does not exist
		dstDir := fmt.Sprintf("%s%s%s", HOME, os.PathSeparator, adler32.Checksum([]byte(absPath)))
		isNew := false

		if _, err = os.Stat(dstDir); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(dstDir, 0750)
				isNew = true
			}
			if err != nil {
				errList = append(errList, err)
				continue
			}
		}

		if isNew {
			pkg, err := ParseDir(absPath)
			if err != nil {
				errList = append(errList, err)
				continue
			}
			if workDir, err := Build(pkg); err != nil {
				errList = append(errList, err)
				continue
			}
		}
	}

	exitCode := 0

	if len(errList) != 0 {
		exitCode = 1
		for _, v := range errList {
			fmt.Fprintf(os.Stderr, "%s", v)
		}
	}

	os.Exit(exitCode)
}
