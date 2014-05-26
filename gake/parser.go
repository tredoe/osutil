// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	IMPORT_PATH = `"github.com/kless/osutil/gake/making"`
	PREFIX      = "Make"
)

// makePackage represents a package of make files.
type makePackage struct {
	Name  string
	Files []makeFile
}

// makeFile represents a set of declarations of make functions.
type makeFile struct {
	Name      string
	MakeFuncs []makeFunc
}

// makeFunc represents a make function.
type makeFunc struct {
	Name string
	Doc  string
}

// The "gake" command expects to find make functions in the "*_make.go" files.
//
// A make function is one named MakeXXX (where XXX is any alphanumeric string
// not starting with a lower case letter) and should have the signature,
//
//	func MakeXXX(m *making.M) { ... }
func ParseDir(path string) (*makePackage, error) {
	filter := func(info os.FileInfo) bool {
		if strings.HasSuffix(info.Name(), "_make.go") {
			return true
		}
		return false
	}

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments|parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, ErrNoMake
	} else if len(pkgs) > 1 {
		return nil, MultiPkgError{path, pkgs}
	}

	pkgName := ""
	for k, _ := range pkgs {
		pkgName = k
		break
	}

	goFiles := make([]makeFile, 0)

	for filename, file := range pkgs[pkgName].Files {
		makeFuncs := make([]makeFunc, 0)

		for _, decl := range file.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			funcName := f.Name.Name

			// Check function name
			if !strings.HasPrefix(funcName, PREFIX) || len(funcName) <= len(PREFIX) {
				continue
			}
			if r, _ := utf8.DecodeRune([]byte(funcName[len(PREFIX):])); !unicode.IsUpper(r) && !unicode.IsDigit(r) {
				continue
			}

			// Check function signature

			if f.Type.Results != nil || len(f.Type.Params.List) != 1 {
				return nil, FuncSignError{fset, file, f}
			}
			pointerType, ok := f.Type.Params.List[0].Type.(*ast.StarExpr)
			if !ok {
				return nil, FuncSignError{fset, file, f}
			}
			selector, ok := pointerType.X.(*ast.SelectorExpr)
			if !ok {
				return nil, FuncSignError{fset, file, f}
			}
			if selector.X.(*ast.Ident).Name != "making" || selector.Sel.Name != "M" {
				return nil, FuncSignError{fset, file, f}
			}

			makeFuncs = append(makeFuncs, makeFunc{funcName, f.Doc.Text()})
		}
		if len(makeFuncs) == 0 {
			continue
		}

		// Check import path
		hasImportPath := false
		for _, v := range file.Imports {
			if v.Path.Value == IMPORT_PATH {
				hasImportPath = true
				break
			}
		}
		if !hasImportPath {
			return nil, ImportPathError{filename}
		}

		// Check the build constraint
		hasBuildCons := false
		for _, c := range file.Comments {
			comment := c.Text()
			if strings.HasPrefix(comment, "+build") {
				words := strings.Split(comment, " ")
				if words[0] == "+build" && words[1] == "gake\n" {
					hasBuildCons = true
					break
				}
			}
		}
		if !hasBuildCons {
			return nil, BuildConsError{filename}
		}

		goFiles = append(goFiles, makeFile{filename, makeFuncs})
	}

	if len(goFiles) == 0 {
		return nil, ErrNoMakeRun
	}
	return &makePackage{pkgName, goFiles}, nil
}

// == Errors
//

var (
	ErrNoMake    = errors.New("  [no make files]")
	ErrNoMakeRun = errors.New("  [no makes to run]")
)

// BuildConsError reports lacking of build constraint.
type BuildConsError struct {
	filename string
}

func (e BuildConsError) Error() string {
	return fmt.Sprintf("%s: no build constraint: \"+build gake\"", e.filename)
}

// FuncSignError represents an incorrect function signature.
type FuncSignError struct {
	fileSet  *token.FileSet
	makeFile *ast.File
	makeFunc *ast.FuncDecl
}

func (e FuncSignError) Error() string {
	return fmt.Sprintf("%s: %s.%s should have the signature func(*making.M)",
		e.fileSet.Position(e.makeFile.Pos()),
		e.makeFile.Name.Name,
		e.makeFunc.Name.Name,
	)
}

// ImportPathError represents a file without a necessary import path.
type ImportPathError struct {
	filename string
}

func (e ImportPathError) Error() string {
	return fmt.Sprintf("%s: no import path: %s", e.filename, IMPORT_PATH)
}

// MultiPkgError represents an error due to multiple packages into a same directory.
type MultiPkgError struct {
	path string
	pkgs map[string]*ast.Package
}

func (e MultiPkgError) Error() string {
	msg := make([]string, len(e.pkgs))
	i := 0

	for pkgName, pkg := range e.pkgs {
		files := make([]string, len(pkg.Files))
		j := 0

		for fileName, _ := range pkg.Files {
			files[j] = "'" + fileName + "'"
			j++
		}

		msg[i] = fmt.Sprintf("%q (%s)", pkgName, strings.Join(files, ", "))
		i++
	}

	return fmt.Sprintf("can't load package: found packages %s in '%s'",
		strings.Join(msg, ", "),
		e.path,
	)
}
