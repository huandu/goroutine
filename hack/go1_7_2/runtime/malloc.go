// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Memory allocator, based on tcmalloc.
// http://goog-perftools.sourceforge.net/doc/tcmalloc.html

// The main allocator works in runs of pages.
// Small allocation sizes (up to and including 32 kB) are
// rounded to one of about 100 size classes, each of which
// has its own free list of objects of exactly that size.
// Any free page of memory can be split into a set of objects
// of one size class, which are then managed using free list
// allocators.
//
// The allocator's data structures are:
//
//	FixAlloc: a free-list allocator for fixed-size objects,
//		used to manage storage used by the allocator.
//	MHeap: the malloc heap, managed at page (4096-byte) granularity.
//	MSpan: a run of pages managed by the MHeap.
//	MCentral: a shared free list for a given size class.
//	MCache: a per-thread (in Go, per-P) cache for small objects.
//	MStats: allocation statistics.
//
// Allocating a small object proceeds up a hierarchy of caches:
//
//	1. Round the size up to one of the small size classes
//	   and look in the corresponding MCache free list.
//	   If the list is not empty, allocate an object from it.
//	   This can all be done without acquiring a lock.
//
//	2. If the MCache free list is empty, replenish it by
//	   taking a bunch of objects from the MCentral free list.
//	   Moving a bunch amortizes the cost of acquiring the MCentral lock.
//
//	3. If the MCentral free list is empty, replenish it by
//	   allocating a run of pages from the MHeap and then
//	   chopping that memory into objects of the given size.
//	   Allocating many objects amortizes the cost of locking
//	   the heap.
//
//	4. If the MHeap is empty or has no page runs large enough,
//	   allocate a new group of pages (at least 1MB) from the
//	   operating system.  Allocating a large run of pages
//	   amortizes the cost of talking to the operating system.
//
// Freeing a small object proceeds up the same hierarchy:
//
//	1. Look up the size class for the object and add it to
//	   the MCache free list.
//
//	2. If the MCache free list is too long or the MCache has
//	   too much memory, return some to the MCentral free lists.
//
//	3. If all the objects in a given span have returned to
//	   the MCentral list, return that span to the page heap.
//
//	4. If the heap has too much memory, return some to the
//	   operating system.
//
//	TODO(rsc): Step 4 is not implemented.
//
// Allocating and freeing a large object uses the page heap
// directly, bypassing the MCache and MCentral free lists.
//
// The small objects on the MCache and MCentral free lists
// may or may not be zeroed. They are zeroed if and only if
// the second word of the object is zero. A span in the
// page heap is zeroed unless s->needzero is set. When a span
// is allocated to break into small objects, it is zeroed if needed
// and s->needzero is set. There are two main benefits to delaying the
// zeroing this way:
//
//	1. stack frames allocated from the small object lists
//	   or the page heap can avoid zeroing altogether.
//	2. the cost of zeroing when reusing a small object is
//	   charged to the mutator, not the garbage collector.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_7_2/runtime/internal/sys"
	"unsafe"
)

const (
	debugMalloc = false

	maxTinySize   = _TinySize
	tinySizeClass = _TinySizeClass
	maxSmallSize  = _MaxSmallSize

	pageShift = _PageShift
	pageSize  = _PageSize
	pageMask  = _PageMask

	maxObjsPerSpan = pageSize / 8

	mSpanInUse = _MSpanInUse

	concurrentSweep = _ConcurrentSweep
)

const (
	_PageShift = 13
	_PageSize  = 1 << _PageShift
	_PageMask  = _PageSize - 1
)

const (
	_64bit = 1 << (^uintptr(0) >> 63) / 2

	_NumSizeClasses = 67

	_MaxSmallSize = 32 << 10

	_TinySize      = 16
	_TinySizeClass = 2

	_FixAllocChunk  = 16 << 10
	_MaxMHeapList   = 1 << (20 - _PageShift)
	_HeapAllocChunk = 1 << 20

	_StackCacheSize = 32 * 1024

	_NumStackOrders = 4 - sys.PtrSize/4*sys.GoosWindows - 1*sys.GoosPlan9

	_MHeapMap_TotalBits = (_64bit*sys.GoosWindows)*35 + (_64bit*(1-sys.GoosWindows)*(1-sys.GoosDarwin*sys.GoarchArm64))*39 + sys.GoosDarwin*sys.GoarchArm64*31 + (1-_64bit)*32
	_MHeapMap_Bits      = _MHeapMap_TotalBits - _PageShift

	_MaxMem = uintptr(1<<_MHeapMap_TotalBits - 1)

	_MaxGcproc = 32
)

const _MaxArena32 = 1<<32 - 1

type persistentAlloc struct {
	base unsafe.Pointer
	off  uintptr
}
