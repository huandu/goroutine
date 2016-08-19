// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Garbage collector: finalizers and block profiling.

package runtime

import "unsafe"

type finblock struct {
	alllink *finblock
	next    *finblock
	cnt     int32
	_       int32
	fin     [(_FinBlockSize - 2*ptrSize - 2*4) / unsafe.Sizeof(finalizer{})]finalizer
}

// NOTE: Layout known to queuefinalizer.
type finalizer struct {
	fn   *funcval
	arg  unsafe.Pointer
	nret uintptr
	fint *_type
	ot   *ptrtype
}
