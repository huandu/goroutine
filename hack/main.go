// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"regexp"

	"github.com/huandu/goroutine/version"
)

var (
	flagGoSrc      string
	flagOutput     string
	flagImportPath string

	reImportPath = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z0-9_]+)*(/[a-zA-Z0-9_]+)*$`)
)

type Context struct {
	GoSrc      string          // Path to go src root directory.
	GoPackage  string          // Path to package src in go src directory.
	Version    version.Version // Version number.
	Output     string          // Path to output directory.
	ImportPath string          // The prefix of import path for output directory.
}

func init() {
	flag.StringVar(&flagGoSrc, "go-src", "", "Path to go src directory.")
	flag.StringVar(&flagOutput, "output", "", "Output path. Default is current directory.")
	flag.StringVar(&flagImportPath, "import-path", "", "The prefix of import path for output directory.")
}

func validateGoSrc(context *Context) {
	if flagGoSrc == "" {
		logFatalf("Flag -go-src must be set.")
	}

	src, err := filepath.Abs(flagGoSrc)

	if err != nil {
		logFatalf("Fail to get absolute path of go src. [go-src:%v]", flagGoSrc)
	}

	info, err := os.Stat(src)

	if err != nil {
		logFatalf("Value of flag -go-src must be a valid directory name. [go-src:%v]", flagGoSrc)
	}

	if !info.Mode().IsDir() {
		logFatalf("Value of flag -go-src must be a directory. [go-src:%v]", flagGoSrc)
	}

	srcsrc := filepath.Join(src, "src")
	info, err = os.Stat(srcsrc)

	if err != nil {
		logFatalf("Fail to find `src` in go src directory. [go-src:%v]", flagGoSrc)
	}

	if !info.Mode().IsDir() {
		logFatalf("The `src` in go src must be a directory. [go-src:%v]", flagGoSrc)
	}

	versionPath := filepath.Join(src, "VERSION")
	versionFile, err := os.Open(versionPath)

	if err != nil {
		logFatalf("Fail to read VERSION in go src directory. Are you sure -go-src points to a valid go src directory? [go-src:%v]", flagGoSrc)
	}

	defer versionFile.Close()
	versionBuf := &bytes.Buffer{}
	versionBuf.ReadFrom(versionFile)

	if versionBuf.Len() == 0 {
		logFatalf("VERSION in go src is empty. [go-src:%v]", flagGoSrc)
	}

	version, err := version.Parse(versionBuf.String())

	if err != nil {
		logFatalf("VERSION file content is not valid. [go-src:%v] [err:%v] [content:%v]", flagGoSrc, err, versionBuf.String())
	}

	context.GoSrc = src
	context.GoPackage = srcsrc
	context.Version = version
}

func validateOutput(context *Context) {
	var err error

	if flagOutput == "" {
		flagOutput, err = os.Getwd()

		if err != nil {
			logFatalf("Value of flag -output is empty and current directory is not available.")
		}
	}

	output, err := filepath.Abs(flagOutput)

	if err != nil {
		logFatalf("Fail to get absolute path of output. [output:%v]", flagOutput)
	}

	info, err := os.Stat(output)

	if err != nil {
		logFatalf("Value of flag -output must be a valid directory name. [output:%v]", flagOutput)
	}

	if !info.Mode().IsDir() {
		logFatalf("Value of flag -output must be a directory. [output:%v]", flagOutput)
	}

	context.Output = output
}

func validateImportPath(context *Context) {
	if flagImportPath == "" {
		logFatalf("Flag -import-path must be set.")
	}

	if !reImportPath.MatchString(flagImportPath) {
		logFatalf("Value of flag -import-path must be a valid import path. [import-path:%v]", flagImportPath)
	}

	context.ImportPath = flagImportPath
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logErrorf("%v", err)
			os.Exit(1)
		}
	}()

	flag.Parse()

	context := &Context{}
	validateGoSrc(context)
	validateOutput(context)
	validateImportPath(context)

	NewGenerator(context, &RuntimeHacker{}).Parse()
}
