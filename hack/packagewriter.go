// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
    "os"
    "path/filepath"
)

var copyright = `// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.
`

type PackageWriter struct {
    Context *Context // Current running context.
    Package string // Package name.
    GoDir string // Hacked go src path name.
    Copyright string // Copyright notice.

    output string // Path for all output files.
}

func NewPackageWriter(context *Context, pkg string) *PackageWriter {
    goDir := "go" + context.Version.Join("_")
	output := filepath.Join(context.Output, goDir, pkg)
    return &PackageWriter{
        Context: context,
        Package: pkg,
        GoDir: goDir,
        Copyright: copyright,

        output: output,
    }
}

func (pw *PackageWriter) Mkdir() error {
    return os.MkdirAll(pw.output, os.ModeDir | 0755)
}

func (pw *PackageWriter) CreateFile(filename string) (*os.File, error) {
    fullPath := filepath.Join(pw.output, filename)
    return os.Create(fullPath)
}
