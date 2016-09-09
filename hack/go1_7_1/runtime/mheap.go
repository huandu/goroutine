// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Page heap.
//
// See malloc.go for overview.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_7_1/runtime/internal/sys"
	"unsafe"
)

// minPhysPageSize is a lower-bound on the physical page size. The
// true physical page size may be larger than this. In contrast,
// sys.PhysPageSize is an upper-bound on the physical page size.
const minPhysPageSize = 4096

// Main malloc heap.
// The heap itself is the "free[]" and "large" arrays,
// but all the other global data is here too.
type mheap struct {
	lock      mutex
	free      [_MaxMHeapList]mSpanList
	freelarge mSpanList
	busy      [_MaxMHeapList]mSpanList
	busylarge mSpanList
	allspans  **mspan
	gcspans   **mspan
	nspan     uint32
	sweepgen  uint32
	sweepdone uint32

	spans        **mspan
	spans_mapped uintptr

	pagesInUse        uint64
	spanBytesAlloc    uint64
	pagesSwept        uint64
	sweepPagesPerByte float64

	largefree  uint64
	nlargefree uint64
	nsmallfree [_NumSizeClasses]uint64

	bitmap         uintptr
	bitmap_mapped  uintptr
	arena_start    uintptr
	arena_used     uintptr
	arena_end      uintptr
	arena_reserved bool

	central [_NumSizeClasses]struct {
		mcentral mcentral
		pad      [sys.CacheLineSize]byte
	}

	spanalloc             fixalloc
	cachealloc            fixalloc
	specialfinalizeralloc fixalloc
	specialprofilealloc   fixalloc
	speciallock           mutex
}

// An MSpan representing actual memory has state _MSpanInUse,
// _MSpanStack, or _MSpanFree. Transitions between these states are
// constrained as follows:
//
// * A span may transition from free to in-use or stack during any GC
//   phase.
//
// * During sweeping (gcphase == _GCoff), a span may transition from
//   in-use to free (as a result of sweeping) or stack to free (as a
//   result of stacks being freed).
//
// * During GC (gcphase != _GCoff), a span *must not* transition from
//   stack or in-use to free. Because concurrent GC may read a pointer
//   and then look up its span, the span state must be monotonic.
const (
	_MSpanInUse = iota
	_MSpanStack
	_MSpanFree
	_MSpanDead
)

// mSpanList heads a linked list of spans.
//
// Linked list structure is based on BSD's "tail queue" data structure.
type mSpanList struct {
	first *mspan
	last  **mspan
}

type mspan struct {
	next *mspan
	prev **mspan
	list *mSpanList

	startAddr     uintptr
	npages        uintptr
	stackfreelist gclinkptr

	freeindex uintptr

	nelems uintptr

	allocCache uint64

	allocBits  *uint8
	gcmarkBits *uint8

	sweepgen    uint32
	divMul      uint32
	allocCount  uint16
	sizeclass   uint8
	incache     bool
	state       uint8
	needzero    uint8
	divShift    uint8
	divShift2   uint8
	elemsize    uintptr
	unusedsince int64
	npreleased  uintptr
	limit       uintptr
	speciallock mutex
	specials    *special
	baseMask    uintptr
}

const (
	_KindSpecialFinalizer = 1
	_KindSpecialProfile   = 2
)

type special struct {
	next   *special
	offset uint16
	kind   byte
}

// The described object has a finalizer set for it.
type specialfinalizer struct {
	special special
	fn      *funcval
	nret    uintptr
	fint    *_type
	ot      *ptrtype
}

// The described object is being heap profiled.
type specialprofile struct {
	special special
	b       *bucket
}

const gcBitsChunkBytes = uintptr(64 << 10)
const gcBitsHeaderBytes = unsafe.Sizeof(gcBitsHeader{})

type gcBitsHeader struct {
	free uintptr
	next uintptr
}

type gcBits struct {
	free uintptr
	next *gcBits
	bits [gcBitsChunkBytes - gcBitsHeaderBytes]uint8
}
