// Copyright 2016 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

type RuntimeHacker struct {}

func (h *RuntimeHacker) Package() string {
    return "runtime"
}

func (h *RuntimeHacker) Hack(pw *PackageWriter) {
    h.genHackedtypes(pw)
}

func (h *RuntimeHacker) genHackedtypes(pw *PackageWriter) {
    file, err := pw.CreateFile("hackedtypes.go")

    if err != nil {
        panic(err)
    }

    defer file.Close()
    file.WriteString(pw.Copyright)
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
