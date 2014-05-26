// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"testing"

	"github.com/kless/goutil"
)

func TestCommand(t *testing.T) {
	commandTests := []goutil.CommandTest{
		{
			Args: "./testdata/",
			//Out:  "",
		},

		{
			Args:   "./testdata/build_cons1/",
			Stderr: BuildConsError{"testdata/build_cons1/test1-constraint_make.go"}.Error() + "\n",
		},
		{
			Args:   "./testdata/build_cons2/",
			Stderr: BuildConsPosError{"testdata/build_cons2/test2-constraint_make.go"}.Error() + "\n",
		},
		{
			Args:   "./testdata/func_sign/",
			Stderr: `testdata/func_sign/test-signature_make.go:3:1: main.MakeTest should have the signature func(*making.M)` + "\n",
		},
		{
			Args:   "./testdata/import_path/",
			Stderr: ImportPathError{"testdata/import_path/test-import_make.go"}.Error() + "\n",
		},
		{
			Args:   "./testdata/multi_pkg/",
			Stderr: `can't load package: found packages "main2" ('testdata/multi_pkg/test3_make.go', 'testdata/multi_pkg/test2_make.go'), "main" ('testdata/multi_pkg/test1_make.go') in './testdata/multi_pkg/'` + "\n",
		},
		{
			Args:   "./testdata/nomake/",
			Stderr: ErrNoMake.Error() + "\n",
		},
		{
			Args:   "./testdata/nomake_run/",
			Stderr: ErrNoMakeRun.Error() + "\n",
		},
	}

	err := goutil.TestingCmd(".", commandTests)
	if err != nil {
		t.Fatal(err)
	}
}
