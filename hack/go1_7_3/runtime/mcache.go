// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// Per-thread (in Go, per-P) cache for small objects.
// No locking needed because it is per-thread (per-P).
//
// mcaches are allocated from non-GC'd memory, so any heap pointers
// must be specially handled.
type mcache struct {
	next_sample int32
	local_scan  uintptr

	tiny             uintptr
	tinyoffset       uintptr
	local_tinyallocs uintptr

	alloc [_NumSizeClasses]*mspan

	stackcache [_NumStackOrders]stackfreelist

	local_nlookup    uintptr
	local_largefree  uintptr
	local_nlargefree uintptr
	local_nsmallfree [_NumSizeClasses]uintptr
}

// A gclink is a node in a linked list of blocks, like mlink,
// but it is opaque to the garbage collector.
// The GC does not trace the pointers during collection,
// and the compiler does not emit write barriers for assignments
// of gclinkptr values. Code should store references to gclinks
// as gclinkptr, not as *gclink.
type gclink struct {
	next gclinkptr
}

// A gclinkptr is a pointer to a gclink, but it is opaque
// to the garbage collector.
type gclinkptr uintptr

type stackfreelist struct {
	list gclinkptr
	size uintptr
}
