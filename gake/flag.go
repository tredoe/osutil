// Copyright 2009 The Go Authors.
// Copyright 2014 Jonas mg
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

/*
// makeFlagSpec defines a flag we know about.
type makeFlagSpec struct {
	name       string
	boolVar    *bool
	passToMake bool // pass to Make
	multiOK    bool // OK to have multiple instances
	present    bool // flag has been seen
}

// makeFlagDefn is the set of flags we process.
var makeFlagDefn = []*makeFlagSpec{
	// local.
	{name: "c", boolVar: makeC},

	// build flags.
	{name: "x", boolVar: makeX},

	// passed to 6.out, adding a "make." prefix to the name if necessary: -v becomes -make.v.
	{name: "cpu", passToMake: true},
	{name: "parallel", passToMake: true},
	{name: "run", passToMake: true},
	{name: "short", boolVar: makeShort, passToMake: true},
	{name: "timeout", passToMake: true},
	{name: "v", boolVar: makeV, passToMake: true},
}

// makeFlags processes the command line, grabbing -x and -c, rewriting known flags
// to have "make" before them, and reading the command line for the 6.out.
// Unfortunately for us, we need to do our own flag processing because "gake"
// grabs some flags but otherwise its command line is just a holding place for
// pkg.making's arguments.
func makeFlags(args []string) (packageNames, passToMake []string) {
	inPkg := false
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			if !inPkg && packageNames == nil {
				// First package name we've seen.
				inPkg = true
			}
			if inPkg {
				packageNames = append(packageNames, args[i])
				continue
			}
		}

		if inPkg {
			// Found an argument beginning with "-"; end of package list.
			inPkg = false
		}

		f, value, extraWord := makeFlag(args, i)
		if f == nil {
			// This is a flag we do not know; we must assume
			// that any args we see after this might be flag
			// arguments, not package names.
			inPkg = false
			if packageNames == nil {
				// make non-nil: we have seen the empty package list
				packageNames = []string{}
			}
			passToMake = append(passToMake, args[i])
			continue
		}
		var err error
		switch f.name {
		// bool flags.
		case "c", "x", "v":
			setBoolFlag(f.boolVar, value)
		case "tags":
			buildContext.BuildTags = strings.Fields(value)
		case "timeout":
			makeTimeout = value
		if extraWord {
			i++
		}
		if f.passToMake {
			passToMake = append(passToMake, "-make."+f.name+"="+value)
		}
	}

	return
}

// makeFlag sees if argument i is a known flag and returns its definition, value, and whether it consumed an extra word.
func makeFlag(args []string, i int) (f *makeFlagSpec, value string, extra bool) {
	arg := args[i]
	if strings.HasPrefix(arg, "--") { // reduce two minuses to one
		arg = arg[1:]
	}
	switch arg {
	case "-?", "-h", "-help":
		usage()
	}
	if arg == "" || arg[0] != '-' {
		return
	}
	name := arg[1:]
	// If there's already "make.", drop it for now.
	name = strings.TrimPrefix(name, "make.")
	equals := strings.Index(name, "=")
	if equals >= 0 {
		value = name[equals+1:]
		name = name[:equals]
	}
	for _, f = range makeFlagDefn {
		if name == f.name {
			// Booleans are special because they have modes -x, -x=true, -x=false.
			if f.boolVar != nil {
				if equals < 0 { // otherwise, it's been set and will be verified in setBoolFlag
					value = "true"
				} else {
					// verify it parses
					setBoolFlag(new(bool), value)
				}
			} else { // Non-booleans must have a value.
				extra = equals < 0
				if extra {
					if i+1 >= len(args) {
						makeSyntaxError("missing argument for flag " + f.name)
					}
					value = args[i+1]
				}
			}
			if f.present && !f.multiOK {
				makeSyntaxError(f.name + " flag may be set only once")
			}
			f.present = true
			return
		}
	}
	f = nil
	return
}

func makeSyntaxError(msg string) {
	fmt.Fprintf(os.Stderr, "gake: %s\n", msg)
	fmt.Fprintf(os.Stderr, "run \"gake help\" for more information\n")
	os.Exit(2)
}
*/
