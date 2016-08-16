// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type visitor func(node ast.Node) bool

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if v(node) {
		return v
	}

	return nil
}

// Generator stores progress of package parser.
type Generator struct {
	parsedPkgs map[string]bool
	pkgs       []string
	context    *Context
}

func (g *Generator) Parse() {
	for i := 0; i < len(g.pkgs); i++ {
		g.parsePkg(g.pkgs[i])
	}
}

func (g *Generator) parsePkg(pkg string) {
	pkgPath := filepath.Join(g.context.GoPackage, pkg)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, func(info os.FileInfo) bool {
		// Filter out all test files.
		if strings.HasSuffix(info.Name(), "_test.go") {
			return false
		}

		return true
	}, parser.ParseComments)

	if err != nil {
		panic(err)
	}

    if _, ok := pkgs["main"]; ok {
        delete(pkgs, "main")
    }

	if len(pkgs) != 1 {
        keys := []string{}

        for k, _ := range pkgs {
            keys = append(keys, k)
        }

		panic(fmt.Errorf("there must be only one package name in a package. [pkgs:%v]", strings.Join(keys, ", ")))
	}

	//goDir := "go" + g.context.Version.Join("_")
	//output := filepath.Join(g.context.Output, goDir, pkg)
	//importPath := filepath.Join(g.context.ImportPath, goDir, pkg)

	for _, p := range pkgs {
		//name := p.Name
		files := p.Files

		for _, f := range files {
			decls := f.Decls
			neededDecls := []ast.Decl{}

			for _, decl := range decls {
				// Only type decl is needed.
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
					neededDecls = append(neededDecls, decl)
					specs := genDecl.Specs

					for _, spec := range specs {
						ast.Walk(visitor(func(node ast.Node) bool {
							if node == nil {
								return false
							}

							switch n := node.(type) {
							case *ast.TypeSpec:
                                typeName := n.Name.Name
								logTracef("Find type. [type:%v]", typeName)

                            case *ast.SelectorExpr:
                                pkgName := n.X.(*ast.Ident).Name
                                typeName := n.Sel.Name

                                if pkgName == "C" {
                                    break
                                }

                                logTracef("Find type. [package-name:%v] [type:%v]", pkgName, typeName)
							}

							return true
						}), spec)
					}
				}
			}

			f.Decls = neededDecls
		}
	}
}

// Generate hacked files for packages.
// Basically, it extracts all types and generates hacked go files.
//
// Panic if it encounters any error.
func GenerateHackedFiles(context *Context, pkgs ...string) {
	generator := &Generator{
		parsedPkgs: make(map[string]bool),
		pkgs:       pkgs,
		context:    context,
	}
	generator.Parse()
}
