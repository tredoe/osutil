// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}

	errList := make([]error, 0)
	workDir := ""

	for _, v := range args {
		pkg, err := ParseDir(v)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		if workDir, err = Build(pkg); err != nil {
			errList = append(errList, err)
			continue
		}
	}

	exitCode := 0

	if len(errList) != 0 {
		exitCode = 1
		for _, v := range errList {
			fmt.Fprintf(os.Stderr, "%s", v)
		}
	}
	/*if err := os.RemoveAll(workDir); err != nil {
		exitCode = 1
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}*/
	println(workDir)

	os.Exit(exitCode)
}
