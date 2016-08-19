// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Page heap.
//
// See malloc.go for overview.

package runtime

// Main malloc heap.
// The heap itself is the "free[]" and "large" arrays,
// but all the other global data is here too.
type mheap struct {
	lock      mutex
	free      [_MaxMHeapList]mspan
	freelarge mspan
	busy      [_MaxMHeapList]mspan
	busylarge mspan
	allspans  **mspan
	gcspans   **mspan
	nspan     uint32
	sweepgen  uint32
	sweepdone uint32

	spans        **mspan
	spans_mapped uintptr

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
		pad      [_CacheLineSize]byte
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
	_MSpanListHead
	_MSpanDead
)

type mspan struct {
	next     *mspan
	prev     *mspan
	start    pageID
	npages   uintptr
	freelist gclinkptr

	sweepgen    uint32
	divMul      uint32
	ref         uint16
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
