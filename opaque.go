// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goroutine

import "unsafe"

const ptrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
const regSize = 4 << (^uintreg(0) >> 63) // unsafe.Sizeof(uintreg(0)) but an ideal const

type mcache struct{}
type typeAlg struct{}
type traceBuf struct{}

type gcMarkWorkerMode int

type persistentAlloc struct {
	base unsafe.Pointer
	off  uintptr
}

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
	// Invariant: wbuf is never full or empty
	wbuf wbufptr

	// Bytes marked (blackened) on this gcWork. This is aggregated
	// into work.bytesMarked by dispose.
	bytesMarked uint64

	// Scan work performed on this gcWork. This is aggregated into
	// gcController by dispose.
	scanWork int64
}

// Information from the compiler about the layout of stack frames.
type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
