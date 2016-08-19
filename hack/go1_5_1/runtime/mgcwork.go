// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

const (
	_Debugwbufs  = false
	_WorkbufSize = 1 * 256
)

// A wbufptr holds a workbuf*, but protects it from write barriers.
// workbufs never live on the heap, so write barriers are unnecessary.
// Write barriers on workbuf pointers may also be dangerous in the GC.
type wbufptr uintptr

// A gcWork provides the interface to produce and consume work for the
// garbage collector.
//
// A gcWork can be used on the stack as follows:
//
//     var gcw gcWork
//     disable preemption
//     .. call gcw.put() to produce and gcw.get() to consume ..
//     gcw.dispose()
//     enable preemption
//
// Or from the per-P gcWork cache:
//
//     (preemption must be disabled)
//     gcw := &getg().m.p.ptr().gcw
//     .. call gcw.put() to produce and gcw.get() to consume ..
//     if gcphase == _GCmarktermination {
//         gcw.dispose()
//     }
//
// It's important that any use of gcWork during the mark phase prevent
// the garbage collector from transitioning to mark termination since
// gcWork may locally hold GC work buffers. This can be done by
// disabling preemption (systemstack or acquirem).
type gcWork struct {
	wbuf wbufptr

	bytesMarked uint64

	scanWork int64
}

type workbufhdr struct {
	node  lfnode
	nobj  int
	inuse bool
	log   [4]int
}

type workbuf struct {
	workbufhdr

	obj [(_WorkbufSize - unsafe.Sizeof(workbufhdr{})) / ptrSize]uintptr
}
