// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Page heap.
//
// See malloc.go for the general overview.
//
// Large spans are the subject of this file. Spans consisting of less than
// _MaxMHeapLists are held in lists of like sized spans. Larger spans
// are held in a treap. See https://en.wikipedia.org/wiki/Treap or
// http://faculty.washington.edu/aragon/pubs/rst89.pdf for an overview.
// sema.go also holds an implementation of a treap.
//
// Each treapNode holds a single span. The treap is sorted by page size
// and for spans of the same size a secondary sort based on start address
// is done.
// Spans are returned based on a best fit algorithm and for spans of the same
// size the one at the lowest address is selected.
//
// The primary routines are
// insert: adds a span to the treap
// remove: removes the span from that treap that best fits the required size
// removeSpan: which removes a specific span from the treap
//
// _mheap.lock must be held when manipulating this data structure.

package runtime

//go:notinheap
type mTreap struct {
	treap *treapNode
}

//go:notinheap
type treapNode struct {
	right     *treapNode
	left      *treapNode
	parent    *treapNode
	npagesKey uintptr
	spanKey   *mspan
	priority  uint32
}
