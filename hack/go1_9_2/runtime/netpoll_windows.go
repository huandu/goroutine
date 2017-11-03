// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const _DWORD_MAX = 0xffffffff

const _INVALID_HANDLE_VALUE = ^uintptr(0)

// net_op must be the same as beginning of internal/poll.operation.
// Keep these in sync.
type net_op struct {
	o overlapped

	pd    *pollDesc
	mode  int32
	errno int32
	qty   uint32
}

type overlappedEntry struct {
	key      uintptr
	op       *net_op
	internal uintptr
	qty      uint32
}
