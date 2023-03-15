// Copyright 2023 Miklos Vajna
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

// HandleOneFile handles AST checks for one file.
func HandleOneFile(filename string) (int, error) {
	parsedAst, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return 0, fmt.Errorf("ParseFile() failed: %s", err)
	}

	code := 0
	for _, decl := range parsedAst.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.CONST {
				continue
			}

			// Flag mutable global variables in random files.
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for _, id := range spec.Names {
						if filename != "commands/context.go" {
							fmt.Printf("%s: %s is a mutable global variable\n", filename, id.Name)
							code = 1
						}
					}
				}
			}
		}
	}

	return code, nil
}

// Main is the commandline interface to this package.
func Main() (int, error) {
	code := 0
	for _, arg := range os.Args[1:] {
		if arg == "--" {
			continue
		}

		ret, err := HandleOneFile(arg)
		if err != nil {
			return 0, fmt.Errorf("HandleOneFile() failed: %s", err)
		}
		code |= ret
	}

	return code, nil
}

func main() {
	code, err := Main()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
