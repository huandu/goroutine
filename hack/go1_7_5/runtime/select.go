// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

const (
	debugSelect = false

	caseRecv = iota
	caseSend
	caseDefault
)

// Select statement header.
// Known to compiler.
// Changes here must also be made in src/cmd/internal/gc/select.go's selecttype.
type hselect struct {
	tcase     uint16
	ncase     uint16
	pollorder *uint16
	lockorder *uint16
	scase     [1]scase
}

// Select case descriptor.
// Known to compiler.
// Changes here must also be made in src/cmd/internal/gc/select.go's selecttype.
type scase struct {
	elem        unsafe.Pointer
	c           *hchan
	pc          uintptr
	kind        uint16
	so          uint16
	receivedp   *bool
	releasetime int64
}

// A runtimeSelect is a single case passed to rselect.
// This must match ../reflect/value.go:/runtimeSelect
type runtimeSelect struct {
	dir selectDir
	typ unsafe.Pointer
	ch  *hchan
	val unsafe.Pointer
}

// These values must match ../reflect/value.go:/SelectDir.
type selectDir int

const (
	_ selectDir = iota
	selectSend
	selectRecv
	selectDefault
)
