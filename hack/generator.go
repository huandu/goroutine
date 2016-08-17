// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/printer"
	"os"
	"path/filepath"
	"strings"
	"regexp"
)

var reBuildIgnoreFlag = regexp.MustCompile(`^/(/|\*)[[:blank:]]*\+build[[:blank:]]+ignore`)

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
	g.parsedPkgs = map[string]bool{}

	for _, pkg := range g.pkgs {
		g.parsedPkgs[pkg] = false
	}

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

	logTracef("Parsed package `%v`...", pkg)

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

	goDir := "go" + g.context.Version.Join("_")
	output := filepath.Join(g.context.Output, goDir, pkg)
	importPrefix := filepath.Join(g.context.ImportPath, goDir, pkg)

	err = os.MkdirAll(output, os.ModeDir | 0755)

	if err != nil {
		panic(err)
	}

	logTracef("Created output path `%v`...", output)

	for _, p := range pkgs {
		//name := p.Name
		files := p.Files

FilesLoop:
		for filename, f := range files {
			filename = filepath.Base(filename)
			logTracef("Package `%v`: Working on file `%v`...", pkg, filename)

			// Skip all ignored files.
			comments := f.Comments

			for _, group := range comments {
				list := group.List

				for _, c := range list {
					if reBuildIgnoreFlag.MatchString(c.Text) {
						logTracef("Package `%v`: Ignored file `%v`.", pkg, filename)
						continue FilesLoop
					}
				}
			}

			imports := f.Imports
			decls := f.Decls
			neededDecls := []ast.Decl{}
			neededImports := []*ast.ImportSpec{}
			checkedImports := map[string]bool{}
			importPathMap := map[string]*ast.ImportSpec{}

			for _, imprt := range imports {
				importName := imprt.Name
				importPath := strings.Trim(imprt.Path.Value, `"`)

				if importName == nil {
					segments := strings.Split(importPath, "/")
					importName = &ast.Ident{
						Name: segments[len(segments) - 1],
					}
				}

				if importName.Name == "." || importName.Name == "_" {
					continue
				}

				importPathMap[importName.Name] = imprt
			}

			for _, decl := range decls {
				genDecl, ok := decl.(*ast.GenDecl)

				// Only type decl is needed.
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				specs := genDecl.Specs
				neededSpecs := []ast.Spec{}

				for _, spec := range specs {
					typeSpec := spec.(*ast.TypeSpec)
					needWalk := true
					ast.Walk(visitor(func(node ast.Node) bool {
						if node == nil {
							return false
						}

						switch n := node.(type) {
                        case *ast.SelectorExpr:
                            pkgName := n.X.(*ast.Ident).Name
                            typeName := n.Sel.Name

							// Ignore entire type spec if it contains any type imported from C.
                            if pkgName == "C" {
								logFatalf("It's not possible to parsee C type `%v.%v` for type. [type:%v]", pkgName, typeName, typeSpec.Name.Name)
								needWalk = false
                                break
                            }

							if _, ok := checkedImports[pkgName]; !ok {
								if importDecl, ok := importPathMap[pkgName]; ok {
									neededImports = append(neededImports, importDecl)
									checkedImports[pkgName] = true
								} else {
									logFatalf("Fail to find package name in import path map. [package:%v] [map:%v]", pkgName, importPathMap)
								}
							}
						}

						return needWalk
					}), spec)

					if needWalk {
						neededSpecs = append(neededSpecs, spec)
					}
				}

				if len(neededSpecs) != 0 {
					genDecl.Specs = neededSpecs
					neededDecls = append(neededDecls, decl)
				}
			}

			// Internal package must be generated if it's referred.
			for _, imprt := range neededImports {
				importPath := strings.Trim(imprt.Path.Value, `"`)
				segments := strings.Split(importPath, "/")

				for _, seg := range segments {
					if seg == "internal" {
						if _, ok := g.parsedPkgs[importPath]; !ok {
							g.pkgs = append(g.pkgs, importPath)
							g.parsedPkgs[importPath] = false
						}

						imprt.Path.Value = `"` + filepath.Join(importPrefix, importPath) + `"`
						break
					}
				}
			}

			if len(neededDecls) == 0 {
				continue
			}

			f.Decls = neededDecls
			f.Imports = neededImports // FIXME: Not work. Hack f.Decls instead. Printer doesn't read it.

			// Create file with modified go src.
			logTracef("Package `%v`: Start to write file `%v`...", pkg, filename)
			fullPath := filepath.Join(output, filename)
			file, err := os.Create(fullPath)

			if err != nil {
				panic(err)
			}

			err = printer.Fprint(file, fset, f)
			file.Close()

			if err != nil {
				panic(err)
			}

			logTracef("Package `%v`: File `%v` is written.", pkg, filename)
		}

		g.parsedPkgs[pkg] = true
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
