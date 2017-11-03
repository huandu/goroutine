// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_9_2/runtime/internal/sys"
	"unsafe"
)

const (
	_WorkbufSize = 2048

	workbufAlloc = 32 << 10
)

// A gcWork provides the interface to produce and consume work for the
// garbage collector.
//
// A gcWork can be used on the stack as follows:
//
//     (preemption must be disabled)
//     gcw := &getg().m.p.ptr().gcw
//     .. call gcw.put() to produce and gcw.get() to consume ..
//     if gcBlackenPromptly {
//         gcw.dispose()
//     }
//
// It's important that any use of gcWork during the mark phase prevent
// the garbage collector from transitioning to mark termination since
// gcWork may locally hold GC work buffers. This can be done by
// disabling preemption (systemstack or acquirem).
type gcWork struct {
	wbuf1, wbuf2 *workbuf

	bytesMarked uint64

	scanWork int64
}

type workbufhdr struct {
	node lfnode
	nobj int
}

//go:notinheap
type workbuf struct {
	workbufhdr

	obj [(_WorkbufSize - unsafe.Sizeof(workbufhdr{})) / sys.PtrSize]uintptr
}
