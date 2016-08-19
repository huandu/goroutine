// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parallel for algorithm.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_6_2/runtime/internal/sys"
)

// A parfor holds state for the parallel for operation.
type parfor struct {
	body   func(*parfor, uint32)
	done   uint32
	nthr   uint32
	thrseq uint32
	cnt    uint32
	wait   bool

	thr []parforthread

	nsteal     uint64
	nstealcnt  uint64
	nprocyield uint64
	nosyield   uint64
	nsleep     uint64
}

// A parforthread holds state for a single thread in the parallel for.
type parforthread struct {
	pos uint64

	nsteal     uint64
	nstealcnt  uint64
	nprocyield uint64
	nosyield   uint64
	nsleep     uint64
	pad        [sys.CacheLineSize]byte
}
