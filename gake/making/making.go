// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// http://golang.org/src/pkg/testing/testing.go

// Package making provides support for automated running of Go packages.
// It is intended to be used in concert with the ``gake'' command, which automates
// execution of any function of the form
//     func MakeXxx(*making.M)
// where Xxx can be any alphanumeric string (but the first letter must not be in
// [a-z]) and serves to identify the run routine.
//
// Within these functions, use the Error, Fail or related methods to signal failure.
//
// To write a new make suite, create a file whose name ends _make.go that
// contains the MakeXxx functions as described here. Put the file in the same
// package as the one being tested. The file will be excluded from regular
// package builds but will be included when the ``gake'' command is run.
// For more detail, run ``gake help''.
//
// Makes may be skipped if not applicable with a call to the Skip method of *M:
//     func MakeTimeConsuming(m *making.M) {
//         if !making.Verbose() {
//             m.Skip("skipping make in non-verbose mode.")
//         }
//         ...
//     }
package making

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// The directory in which to create files and the like. When run from
	// "gake", the binary always runs in the source directory for the package;
	// this flag lets "gake" tell the binary to write the files in the directory
	//  where the "gake" command is run.
	outputDir = flag.String("make.outputdir", "", "directory in which to write files")

	chatty     = flag.Bool("make.v", false, "verbose: print additional output")
	match      = flag.String("make.run", "", "regular expression to select runs")
	timeout    = flag.Duration("make.timeout", 0, "if positive, sets an aggregate time limit for all runs")
	cpuListStr = flag.String("make.cpu", "", "comma-separated list of number of CPUs to use for each run")
	parallel   = flag.Int("make.parallel", runtime.GOMAXPROCS(0), "maximum run parallelism")

	cpuList []int
)

// common holds the elements common for R and
// captures common methods such as Errorf.
type common struct {
	mu       sync.RWMutex // guards output and failed
	output   []byte       // Output generated.
	failed   bool         // Make has failed.
	skipped  bool         // Make has been skipped.
	finished bool

	start    time.Time // Time run started
	duration time.Duration
	self     interface{}      // To be sent on signal channel when done.
	signal   chan interface{} // Output for serial executions.
}

// Verbose reports whether the -make.v flag is set.
func Verbose() bool {
	return *chatty
}

func (c *common) private() {}

// Fail marks the function as having failed but continues execution.
func (c *common) Fail() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failed = true
}

// Failed reports whether the function has failed.
func (c *common) Failed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.failed
}

// FailNow marks the function as having failed and stops its execution.
// Execution will continue at the next run.
// FailNow must be called from the goroutine running the
// run function, not from other goroutines
// created during the run. Calling FailNow does not stop
// those other goroutines.
func (c *common) FailNow() {
	c.Fail()

	c.finished = true
	runtime.Goexit()
}

// log generates the output. It's always at the same stack depth.
func (c *common) log(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.output = append(c.output, decorate(s)...)
}

// Log formats its arguments using default formatting, analogous to Println,
// and records the text in the error log. The text will be printed only if
// the run fails or the -make.v flag is set.
func (c *common) Log(args ...interface{}) { c.log(fmt.Sprintln(args...)) }

// Logf formats its arguments according to the format, analogous to Printf,
// and records the text in the error log. The text will be printed only if
// the run fails or the -make.v flag is set.
func (c *common) Logf(format string, args ...interface{}) { c.log(fmt.Sprintf(format, args...)) }

// Error is equivalent to Log followed by Fail.
func (c *common) Error(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.Fail()
}

// Errorf is equivalent to Logf followed by Fail.
func (c *common) Errorf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.Fail()
}

// Fatal is equivalent to Log followed by FailNow.
func (c *common) Fatal(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.FailNow()
}

// Fatalf is equivalent to Logf followed by FailNow.
func (c *common) Fatalf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.FailNow()
}

// Skip is equivalent to Log followed by SkipNow.
func (c *common) Skip(args ...interface{}) {
	c.log(fmt.Sprintln(args...))
	c.SkipNow()
}

// Skipf is equivalent to Logf followed by SkipNow.
func (c *common) Skipf(format string, args ...interface{}) {
	c.log(fmt.Sprintf(format, args...))
	c.SkipNow()
}

// SkipNow marks the run as having been skipped and stops its execution.
// Execution will continue at the next run. See also FailNow.
// SkipNow must be called from the goroutine running the run, not from
// other goroutines created during the run. Calling SkipNow does not stop
// those other goroutines.
func (c *common) SkipNow() {
	c.skip()
	c.finished = true
	runtime.Goexit()
}

func (c *common) skip() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.skipped = true
}

// Skipped reports whether the run was skipped.
func (c *common) Skipped() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.skipped
}

// * * *

// M is a type passed to Make functions to manage run state and support formatted run logs.
// Logs are accumulated during execution and dumped to standard error when done.
type M struct {
	common
	name          string    // Name of run.
	startParallel chan bool // Parallel tests will wait on this.
}

// Parallel signals that this run is to be run in parallel with (and only with)
// other parallel runs.
func (m *M) Parallel() {
	m.signal <- (*M)(nil) // Release main running loop
	<-m.startParallel     // Wait for serial runs to finish
	// Assuming Parallel is the first thing a run does, which is reasonable,
	// reinitialize the run's start time because it's actually starting now.
	m.start = time.Now()
}

// * * *

// decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func decorate(s string) string {
	_, file, line, ok := runtime.Caller(3) // decorate + log + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	fmt.Fprintf(buf, "%s:%d: ", file, line)
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		buf.WriteString(line)
	}
	buf.WriteByte('\n')
	return buf.String()
}

// toOutputDir returns the file name relocated, if required, to outputDir.
// Simple implementation to avoid pulling in path/filepath.
func toOutputDir(path string) string {
	if *outputDir == "" || path == "" {
		return path
	}
	if runtime.GOOS == "windows" {
		// On Windows, it's clumsy, but we can be almost always correct
		// by just looking for a drive letter and a colon.
		// Absolute paths always have a drive letter (ignoring UNC).
		// Problem: if path == "C:A" and outputdir == "C:\Go" it's unclear
		// what to do, but even then path/filepath doesn't help.
		// TODO: Worth doing better? Probably not, because we're here only
		// under the management of go test.
		if len(path) >= 2 {
			letter, colon := path[0], path[1]
			if ('a' <= letter && letter <= 'z' || 'A' <= letter && letter <= 'Z') && colon == ':' {
				// If path starts with a drive letter we're stuck with it regardless.
				return path
			}
		}
	}
	if os.IsPathSeparator(path[0]) {
		return path
	}
	return fmt.Sprintf("%s%c%s", *outputDir, os.PathSeparator, path)
}

var timer *time.Timer

// startAlarm starts an alarm if requested.
func startAlarm() {
	if *timeout > 0 {
		timer = time.AfterFunc(*timeout, func() {
			panic(fmt.Sprintf("run timed out after %v", *timeout))
		})
	}
}

// stopAlarm turns off the alarm.
func stopAlarm() {
	if *timeout > 0 {
		timer.Stop()
	}
}

func parseCpuList() {
	for _, val := range strings.Split(*cpuListStr, ",") {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		cpu, err := strconv.Atoi(val)
		if err != nil || cpu <= 0 {
			fmt.Fprintf(os.Stderr, "making: invalid value %q for -make.cpu\n", val)
			os.Exit(1)
		}
		cpuList = append(cpuList, cpu)
	}
	if cpuList == nil {
		cpuList = append(cpuList, runtime.GOMAXPROCS(-1))
	}
}
