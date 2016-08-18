// Copyright 2016 Huan Du. All rights reserved.
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
	context    *Context
	parsedPkgs map[string]bool
	pkgs       []string
	hackers map[string]Hacker
}

func NewGenerator(context *Context, hackers ...Hacker) *Generator {
	parsedPkgs := map[string]bool{}
	pkgs := []string{}
	hackersMap := map[string]Hacker{}

	for _, hacker := range hackers {
		pkg := hacker.Package()
		parsedPkgs[pkg] = false
		pkgs = append(pkgs, pkg)
		hackersMap[pkg] = hacker
	}

	return &Generator{
		context:    context,
		parsedPkgs: parsedPkgs,
		pkgs:       pkgs,
		hackers: hackersMap,
	}
}

// Generate hacked files for packages.
// Basically, it extracts all types and generates hacked go files.
//
// Panic if it encounters any error.
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

	logDebugf("Parsed package `%v`...", pkg)

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

	packageWriter := NewPackageWriter(g.context, pkg)
	err = packageWriter.Mkdir()

	if err != nil {
		panic(err)
	}

	var p *ast.Package

	for _, p = range pkgs {
		break
	}

	files := p.Files
	importPrefix := filepath.Join(g.context.ImportPath, packageWriter.GoDir)

FilesLoop:
	for filename, f := range files {
		filename = filepath.Base(filename)
		logDebugf("Package `%v`: Working on file `%v`...", pkg, filename)

		// Skip all ignored files.
		comments := f.Comments

		for _, group := range comments {
			list := group.List

			for _, c := range list {
				if reBuildIgnoreFlag.MatchString(c.Text) {
					logDebugf("Package `%v`: Ignored file `%v`.", pkg, filename)
					continue FilesLoop
				}
			}
		}

		imports := f.Imports
		decls := f.Decls
		neededDecls := []ast.Decl{}
		neededComments := map[*ast.CommentGroup]bool{}
		importPathMap := map[string]string{}
		usedImports := map[string]bool{}
		usedImportPathMap := map[string]bool{}
		allImportDecls := []*ast.GenDecl{}

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

			importPathMap[importName.Name] = importPath
		}

		for _, decl := range decls {
			genDecl, ok := decl.(*ast.GenDecl)

			if !ok {
				continue
			}

			if genDecl.Tok == token.IMPORT {
				allImportDecls = append(allImportDecls, genDecl)
				continue
			}

			if genDecl.Tok != token.TYPE && genDecl.Tok != token.CONST {
				continue
			}

			specs := genDecl.Specs
			neededSpecs := []ast.Spec{}

			for _, spec := range specs {
				needWalk := true
				ast.Walk(visitor(func(node ast.Node) bool {
					if node == nil {
						return false
					}

					switch n := node.(type) {
                    case *ast.SelectorExpr:
						if _, ok := n.X.(*ast.Ident); !ok {
							break
						}

                        pkgName := n.X.(*ast.Ident).Name
                        typeName := n.Sel.Name

						// Ignore entire type spec if it contains any type imported from C.
                        if pkgName == "C" {
							logErrorf("It's not possible to parsee C type `%v.%v` for type. [type:%v]", pkgName, typeName, spec)
							needWalk = false
                            break
                        }

						if _, ok := usedImports[pkgName]; !ok {
							importPath, ok := importPathMap[pkgName]

							if !ok {
								logErrorf("Fail to find package name in import path map. [package:%v] [map:%v]", pkgName, importPathMap)
								needWalk = false
								break
							}

							usedImports[pkgName] = true
							usedImportPathMap[importPath] = true
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

				if genDecl.Doc != nil {
					neededComments[genDecl.Doc] = true
				}
			}
		}

		if len(neededDecls) == 0 {
			logDebugf("Package `%v`: Skip file `%v` as there is no type decl in this file.", pkg, filename)
			continue
		}

		hackedDecls := make([]ast.Decl, 0, len(neededDecls) + len(allImportDecls))
		hackedImportSpecs := make([]*ast.ImportSpec, 0, len(allImportDecls))
		hackedComments := make([]*ast.CommentGroup, 0, len(f.Comments))

		// Filter import decls and hack internal package.
		for _, genDecl := range allImportDecls {
			specs := genDecl.Specs
			neededSpecs := []ast.Spec{}

			for _, spec := range specs {
				importSpec := spec.(*ast.ImportSpec)
				importPath := strings.Trim(importSpec.Path.Value, `"`)

				if _, ok := usedImportPathMap[importPath]; !ok {
					continue
				}

				segments := strings.Split(importPath, "/")

				// Internal package must be generated to make compiler happy.
				for _, seg := range segments {
					if seg == "internal" {
						if _, ok := g.parsedPkgs[importPath]; !ok {
							g.pkgs = append(g.pkgs, importPath)
							g.parsedPkgs[importPath] = false
						}

						importSpec.Path.Value = `"` + filepath.Join(importPrefix, importPath) + `"`
						break
					}
				}

				hackedImportSpecs = append(hackedImportSpecs, importSpec)
				neededSpecs = append(neededSpecs, importSpec)
			}

			if len(neededSpecs) == 0 {
				continue
			}

			genDecl.Specs = neededSpecs
			hackedDecls = append(hackedDecls, genDecl)
		}

		hackedDecls = append(hackedDecls, neededDecls...)

		// Keep comments before `package` and type decl.
		for _, commentGroup := range f.Comments {
			if _, ok := neededComments[commentGroup]; ok || commentGroup.List[0].Slash < f.Package {
				hackedComments = append(hackedComments, commentGroup)
			}
		}

		f.Decls = hackedDecls
		f.Imports = hackedImportSpecs
		f.Comments = hackedComments

		// Create file with modified go src.
		logDebugf("Package `%v`: Start to write file `%v`...", pkg, filename)
		file, err := packageWriter.CreateFile(filename)

		if err != nil {
			panic(err)
		}

		err = printer.Fprint(file, fset, f)
		file.Close()

		if err != nil {
			panic(err)
		}

		logDebugf("Package `%v`: File `%v` is written.", pkg, filename)
	}

	g.parsedPkgs[pkg] = true

	if hacker, ok := g.hackers[pkg]; ok {
		hacker.Hack(packageWriter)
	}
}
