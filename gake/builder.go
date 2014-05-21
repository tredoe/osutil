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

	dstFile := "foo"
	if runtime.GOOS == "windows" {
		dstFile += ".exe"
	}
	cmd := exec.Command("go", "build", "--tags", "gake", "-o", dstFile)
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}

	return
}

var testmainTmpl = template.Must(template.New("main").Parse(`
package main

import (
	"regexp"
	"testing"

	"github.com/kless/osutil/gake/making"

{{if .NeedTest}}
	_test {{.Package.ImportPath | printf "%q"}}
{{end}}
)

var tests = []testing.InternalTest{
{{range .Tests}}
	{"{{.Name}}", {{.Package}}.{{.Name}}},
{{end}}
}

var examples = []testing.InternalExample{
{{range .Examples}}
	{"{{.Name}}", {{.Package}}.{{.Name}}, {{.Output | printf "%q"}}},
{{end}}
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
	testing.Main(matchString, tests, benchmarks, examples)
}

`))
