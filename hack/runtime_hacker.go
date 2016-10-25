// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/huandu/goroutine/copyright"
	"github.com/huandu/goroutine/version"
)

type RuntimeHacker struct{}

func (h *RuntimeHacker) Package() string {
	return "runtime"
}

func (h *RuntimeHacker) Hack(pw *PackageWriter) {
	h.genHackedtypes(pw)

	// This hack only applies to go1.7.2 or later.
	if pw.Context.Version.Compare(version.Version{"1", "7", "2"}) >= 0 {
		h.genHackedBuildruntime(pw)
	}
}

func (h *RuntimeHacker) genHackedtypes(pw *PackageWriter) {
	file := pw.MustCreateFile("hackedtypes.go")
	defer file.Close()
	file.WriteString(copyright.COPYRIGHT)
	file.WriteString(`

package runtime

// Goroutine is the internal type represents a goroutine.
type Goroutine g

// Get goid.
func (g *Goroutine) Goid() int64 {
    return g.goid
}
`)
}

func (h *RuntimeHacker) genHackedBuildruntime(pw *PackageWriter) {
	file := pw.MustCreateFile("internal/sys/hackedbuildruntime.go")
	defer file.Close()
	file.WriteString(copyright.COPYRIGHT)
	file.WriteString(`

package sys

`)
	//file.WriteString(fmt.Sprintf("const DefaultGoroot = `%s`\n", runtime.GOROOT()))
	file.WriteString(fmt.Sprintf("const TheVersion = `%s`\n", pw.Context.Version))
	//file.WriteString(fmt.Sprintf("const Goexperiment = `%s`\n", ))

	// HACK(huandu): Always assume current build is optimized version.
	// If code is built with `-gcflags -N`, this hack fails.
	file.WriteString(fmt.Sprintf("const StackGuardMultiplier = %d\n", 1))
}
