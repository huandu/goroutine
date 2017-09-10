// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Page heap.
//
// See malloc.go for overview.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_9/runtime/internal/sys"
	"unsafe"
)

// minPhysPageSize is a lower-bound on the physical page size. The
// true physical page size may be larger than this. In contrast,
// sys.PhysPageSize is an upper-bound on the physical page size.
const minPhysPageSize = 4096

// Main malloc heap.
// The heap itself is the "free[]" and "large" arrays,
// but all the other global data is here too.
//
// mheap must not be heap-allocated because it contains mSpanLists,
// which must not be heap-allocated.
//
//go:notinheap
type mheap struct {
	lock      mutex
	free      [_MaxMHeapList]mSpanList
	freelarge mTreap
	busy      [_MaxMHeapList]mSpanList
	busylarge mSpanList
	sweepgen  uint32
	sweepdone uint32
	sweepers  uint32

	allspans []*mspan

	spans []*mspan

	sweepSpans [2]gcSweepBuf

	_ uint32

	pagesInUse         uint64
	pagesSwept         uint64
	pagesSweptBasis    uint64
	sweepHeapLiveBasis uint64
	sweepPagesPerByte  float64

	largealloc  uint64
	nlargealloc uint64
	largefree   uint64
	nlargefree  uint64
	nsmallfree  [_NumSizeClasses]uint64

	bitmap        uintptr
	bitmap_mapped uintptr

	arena_start uintptr
	arena_used  uintptr

	arena_alloc uintptr
	arena_end   uintptr

	arena_reserved bool

	_ uint32

	central [numSpanClasses]struct {
		mcentral mcentral
		pad      [sys.CacheLineSize - unsafe.Sizeof(mcentral{})%sys.CacheLineSize]byte
	}

	spanalloc             fixalloc
	cachealloc            fixalloc
	treapalloc            fixalloc
	specialfinalizeralloc fixalloc
	specialprofilealloc   fixalloc
	speciallock           mutex
}

// An MSpan representing actual memory has state _MSpanInUse,
// _MSpanManual, or _MSpanFree. Transitions between these states are
// constrained as follows:
//
// * A span may transition from free to in-use or manual during any GC
//   phase.
//
// * During sweeping (gcphase == _GCoff), a span may transition from
//   in-use to free (as a result of sweeping) or manual to free (as a
//   result of stacks being freed).
//
// * During GC (gcphase != _GCoff), a span *must not* transition from
//   manual or in-use to free. Because concurrent GC may read a pointer
//   and then look up its span, the span state must be monotonic.
type mSpanState uint8

const (
	_MSpanDead mSpanState = iota
	_MSpanInUse
	_MSpanManual
	_MSpanFree
)

// mSpanList heads a linked list of spans.
//
//go:notinheap
type mSpanList struct {
	first *mspan
	last  *mspan
}

//go:notinheap
type mspan struct {
	next *mspan
	prev *mspan
	list *mSpanList

	startAddr uintptr
	npages    uintptr

	manualFreeList gclinkptr

	freeindex uintptr

	nelems uintptr

	allocCache uint64

	allocBits  *gcBits
	gcmarkBits *gcBits

	sweepgen    uint32
	divMul      uint16
	baseMask    uint16
	allocCount  uint16
	spanclass   spanClass
	incache     bool
	state       mSpanState
	needzero    uint8
	divShift    uint8
	divShift2   uint8
	elemsize    uintptr
	unusedsince int64
	npreleased  uintptr
	limit       uintptr
	speciallock mutex
	specials    *special
}

// A spanClass represents the size class and noscan-ness of a span.
//
// Each size class has a noscan spanClass and a scan spanClass. The
// noscan spanClass contains only noscan objects, which do not contain
// pointers and thus do not need to be scanned by the garbage
// collector.
type spanClass uint8

const (
	numSpanClasses = _NumSizeClasses << 1
	tinySpanClass  = spanClass(tinySizeClass<<1 | 1)
)

const (
	_KindSpecialFinalizer = 1
	_KindSpecialProfile   = 2
)

//go:notinheap
type special struct {
	next   *special
	offset uint16
	kind   byte
}

// The described object has a finalizer set for it.
//
// specialfinalizer is allocated from non-GC'd memory, so any heap
// pointers must be specially handled.
//
//go:notinheap
type specialfinalizer struct {
	special special
	fn      *funcval
	nret    uintptr
	fint    *_type
	ot      *ptrtype
}

// The described object is being heap profiled.
//
//go:notinheap
type specialprofile struct {
	special special
	b       *bucket
}

// gcBits is an alloc/mark bitmap. This is always used as *gcBits.
//
//go:notinheap
type gcBits uint8

const gcBitsChunkBytes = uintptr(64 << 10)
const gcBitsHeaderBytes = unsafe.Sizeof(gcBitsHeader{})

type gcBitsHeader struct {
	free uintptr
	next uintptr
}

//go:notinheap
type gcBitsArena struct {
	free uintptr
	next *gcBitsArena
	bits [gcBitsChunkBytes - gcBitsHeaderBytes]gcBits
}
