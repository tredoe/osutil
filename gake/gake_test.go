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
			Out:  "",
		},
		{
			Args:   "./testdata/multi_pkg/",
			Stderr: `can't load package: found packages "main2" ('testdata/multi_pkg/test3_make.go', 'testdata/multi_pkg/test2_make.go'), "main" ('testdata/multi_pkg/test1_make.go') in './testdata/multi_pkg/'` + "\n",
			//Out:  "",
		},
	}

	err := goutil.TestingCmd(".", commandTests)
	if err != nil {
		t.Fatal(err)
	}
}
