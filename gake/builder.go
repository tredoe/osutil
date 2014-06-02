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
	"runtime"
	"text/template"
)

// Build uses the tool "go build" to compile the make files.
// Returns the working directory and the error, if any.
func Build(pkg *makePackage) error {
	workDir, err := ioutil.TempDir("", "gake-")
	if err != nil {
		return
	}
	defer os.RemoveAll(workDir)

	// Copy all files to the temporary directory.
	for _, f := range pkg.Files {
		src, err := ioutil.ReadFile(f.Name)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(workDir+string(os.PathSeparator)+filepath.Base(f.Name), src, 0644)
		if err != nil {
			return err
		}
	}

	// Write the 'makemain.go' file.
	f, err := os.Create(workDir + string(os.PathSeparator) + "makemain.go")
	if err != nil {
		return err
	}
	defer f.Close()
	if err = makemainTmpl.Execute(f, pkg); err != nil {
		return err
	}

	// Build

	cmdName := workDir + string(os.PathSeparator) + CMD_NAME
	if runtime.GOOS == "windows" {
		cmdName += ".exe"
	}
	cmd := new(exec.Cmd)

	if !*makeX {
		cmd = exec.Command("go", "build", "--tags", "gake", "-o", cmdName)
	} else {
		cmd = exec.Command("go", "build", "--tags", "gake", "-o", cmdName, "-x")
	}
	cmd.Dir = workDir
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	// Run
	if !*makeC {
		cmd = exec.Command(cmdName, getMakeArgs()...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	return
}

var makemainTmpl = template.Must(template.New("main").Parse(`
package main

import (
	"regexp"

	"github.com/kless/osutil/gake/making"
)

var makes = []making.InternalMake{
{{range $_, $f := .Files}}{{range $f.MakeFuncs}}
	{"{{.Name}}", {{.Name}}},{{end}}{{end}}
}

var matchPat string
var matchRe *regexp.Regexp

func matchString(pat, str string) (result bool, err error) {
	if matchRe == nil || matchPat != pat {
		matchPat = pat
		matchRe, err = regexp.Compile(matchPat)
		if err != nil {
			return
		}
	}
	return matchRe.MatchString(str), nil
}

func main() {
	making.Main(matchString, makes)
}
`))
