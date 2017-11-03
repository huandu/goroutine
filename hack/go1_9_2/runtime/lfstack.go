// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Lock-free stack.

package runtime

// lfstack is the head of a lock-free stack.
//
// The zero value of lfstack is an empty list.
//
// This stack is intrusive. Nodes must embed lfnode as the first field.
//
// The stack does not keep GC-visible pointers to nodes, so the caller
// is responsible for ensuring the nodes are not garbage collected
// (typically by allocating them from manually-managed memory).
type lfstack uint64
