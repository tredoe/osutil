// Copyright 2014 Jonas mg
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const PREFIX = "Make"

type funcDecl struct {
	name string
	doc  string
}

// The "gake" command expects to find make functions in the "*_make.go" files.
//
// A make function is one named MakeXXX (where XXX is any alphanumeric string
// not starting with a lower case letter) and should have the signature,
//
//	func MakeXXX(m *making.M) { ... }
func Parser(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*_make.go"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Printf("  [no make files]\n")
		return nil
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, files[0], nil, parser.ParseComments)
	if err != nil {
		return err
	}

	funcDecls := make([]funcDecl, 0)

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
			return FuncSignError{fset, file, f}
		}
		pointerType, ok := f.Type.Params.List[0].Type.(*ast.StarExpr)
		if !ok {
			return FuncSignError{fset, file, f}
		}
		selector, ok := pointerType.X.(*ast.SelectorExpr)
		if !ok {
			return FuncSignError{fset, file, f}
		}
		if selector.X.(*ast.Ident).Name != "making" || selector.Sel.Name != "M" {
			return FuncSignError{fset, file, f}
		}

		funcDecls = append(funcDecls, funcDecl{funcName, f.Doc.Text()})
	}

	if len(funcDecls) == 0 {
		fmt.Printf("  [no makes to run]\n")
	}
	return nil
}

// * * *

// NoShellError represents an incorrect function signature.
type FuncSignError struct {
	fileSet  *token.FileSet
	file     *ast.File
	funcDecl *ast.FuncDecl
}

func (e FuncSignError) Error() string {
	return fmt.Sprintf("%s: %s.%s should have the signature func(*making.M)",
		e.fileSet.Position(e.file.Pos()),
		e.file.Name.Name,
		e.funcDecl.Name.Name,
	)
}
