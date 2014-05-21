// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"log"

	"github.com/kless/goutil"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	} else if len(args) > 1 {
		// TODO: error
	}

	pkg, err := ParseDir(args[0])
	if err != nil {
		goutil.Fatalf("%s", err)
	}

	if err = Builder(pkg); err != nil {
		goutil.Fatalf("%s", err)
	}
}
