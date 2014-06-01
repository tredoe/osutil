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
	"time"
)

var makeUsage = func() {
	fmt.Fprintf(os.Stderr, `Usage of gake:

  -c=false: compile but do not run the binary
  -x=false: print command lines as they are executed

  // These flags (used by gake/making) can be passed with or without a "make."
  // prefix: -v or -make.v
  -cpu="": passes -make.cpu
  -parallel=0: passes -make.parallel
  -run="": passes -make.run
  -short=false: passes -make.short
  -timeout=0: passes -make.timeout
  -v=false: passes -make.v
`)
	os.Exit(2)
}

var (
	makeC = flag.Bool("c", false, "compile but do not run the binary")
	makeX = flag.Bool("x", false, "print command lines as they are executed")

	makeCPU      string
	makeParallel int
	makeRun      string
	makeShort    bool
	makeTimeout  time.Duration
	makeV        bool
)

func init() {
	flag.StringVar(&makeCPU, "cpu", "", "passes -make.cpu")
	flag.StringVar(&makeCPU, "make.cpu", "", "")

	flag.IntVar(&makeParallel, "parallel", 0, "passes -make.parallel")
	flag.IntVar(&makeParallel, "make.parallel", 0, "")

	flag.StringVar(&makeRun, "run", "", "passes -make.run")
	flag.StringVar(&makeRun, "make.run", "", "")

	flag.BoolVar(&makeShort, "short", false, "passes -make.short")
	flag.BoolVar(&makeShort, "make.short", false, "")

	flag.DurationVar(&makeTimeout, "timeout", 0, "passes -make.timeout")
	flag.DurationVar(&makeTimeout, "make.timeout", 0, "")

	flag.BoolVar(&makeV, "v", false, "passes -make.v")
	flag.BoolVar(&makeV, "make.v", false, "")

	flag.Usage = makeUsage
}

var (
	makeNeedBinary   bool // need to keep binary around
	makeShowPass     bool // show passing output
	makeStreamOutput bool // show output as it is generated

	makeKillTimeout = 3 * time.Minute
)

// makeArgs returns the arguments to be passed to "making".
func makeArgs() []string {
	args := make([]string, 0)

	// Rewrite known flags to have "make" before them
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "cpu", "parallel", "run", "short", "timeout", "v":
			f.Name = "make." + f.Name
		}

		args = append(args, "-"+f.Name)
		args = append(args, f.Value.String())
	})

	return args
}
