// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Garbage collector: finalizers and block profiling.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_9/runtime/internal/sys"
	"unsafe"
)

// finblock is an array of finalizers to be executed. finblocks are
// arranged in a linked list for the finalizer queue.
//
// finblock is allocated from non-GC'd memory, so any heap pointers
// must be specially handled. GC currently assumes that the finalizer
// queue does not grow during marking (but it can shrink).
//
//go:notinheap
type finblock struct {
	alllink *finblock
	next    *finblock
	cnt     uint32
	_       int32
	fin     [(_FinBlockSize - 2*sys.PtrSize - 2*4) / unsafe.Sizeof(finalizer{})]finalizer
}

// NOTE: Layout known to queuefinalizer.
type finalizer struct {
	fn   *funcval
	arg  unsafe.Pointer
	nret uintptr
	fint *_type
	ot   *ptrtype
}
