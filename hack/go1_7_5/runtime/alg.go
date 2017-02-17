// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_7_5/runtime/internal/sys"
	"unsafe"
)

const (
	c0 = uintptr((8-sys.PtrSize)/4*2860486313 + (sys.PtrSize-4)/4*33054211828000289)
	c1 = uintptr((8-sys.PtrSize)/4*3267000013 + (sys.PtrSize-4)/4*23344194077549503)
)

// type algorithms - known to compiler
const (
	alg_NOEQ = iota
	alg_MEM0
	alg_MEM8
	alg_MEM16
	alg_MEM32
	alg_MEM64
	alg_MEM128
	alg_STRING
	alg_INTER
	alg_NILINTER
	alg_FLOAT32
	alg_FLOAT64
	alg_CPLX64
	alg_CPLX128
	alg_max
)

// typeAlg is also copied/used in reflect/type.go.
// keep them in sync.
type typeAlg struct {
	hash func(unsafe.Pointer, uintptr) uintptr

	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

const hashRandomBytes = sys.PtrSize / 4 * 64
