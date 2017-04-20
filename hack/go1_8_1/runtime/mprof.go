// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Malloc profiling.
// Patterned after tcmalloc's algorithms; shorter code.

package runtime

const (
	memProfile bucketType = 1 + iota
	blockProfile
	mutexProfile

	buckHashSize = 179999

	maxStack = 32
)

type bucketType int

// A bucket holds per-call-stack profiling information.
// The representation is a bit sleazy, inherited from C.
// This struct defines the bucket header. It is followed in
// memory by the stack words and then the actual record
// data, either a memRecord or a blockRecord.
//
// Per-call-stack profiling information.
// Lookup by hashing call stack into a linked-list hash table.
//
// No heap pointers.
//
//go:notinheap
type bucket struct {
	next    *bucket
	allnext *bucket
	typ     bucketType
	hash    uintptr
	size    uintptr
	nstk    uintptr
}

// A memRecord is the bucket data for a bucket of type memProfile,
// part of the memory profile.
type memRecord struct {
	allocs      uintptr
	frees       uintptr
	alloc_bytes uintptr
	free_bytes  uintptr

	prev_allocs      uintptr
	prev_frees       uintptr
	prev_alloc_bytes uintptr
	prev_free_bytes  uintptr

	recent_allocs      uintptr
	recent_frees       uintptr
	recent_alloc_bytes uintptr
	recent_free_bytes  uintptr
}

// A blockRecord is the bucket data for a bucket of type blockProfile,
// which is used in blocking and mutex profiles.
type blockRecord struct {
	count  int64
	cycles int64
}

// A StackRecord describes a single execution stack.
type StackRecord struct {
	Stack0 [32]uintptr
}

// A MemProfileRecord describes the live objects allocated
// by a particular call sequence (stack trace).
type MemProfileRecord struct {
	AllocBytes, FreeBytes     int64
	AllocObjects, FreeObjects int64
	Stack0                    [32]uintptr
}

// BlockProfileRecord describes blocking events originated
// at a particular call sequence (stack trace).
type BlockProfileRecord struct {
	Count  int64
	Cycles int64
	StackRecord
}
