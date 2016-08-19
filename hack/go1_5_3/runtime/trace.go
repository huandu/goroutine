// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go execution tracer.
// The tracer captures a wide range of execution events like goroutine
// creation/blocking/unblocking, syscall enter/exit/block, GC-related events,
// changes of heap size, processor start/stop, etc and writes them to a buffer
// in a compact form. A precise nanosecond-precision timestamp and a stack
// trace is captured for most events.
// See https://golang.org/s/go15trace for more info.

package runtime

import "unsafe"

// Event types in the trace, args are given in square brackets.
const (
	traceEvNone           = 0
	traceEvBatch          = 1
	traceEvFrequency      = 2
	traceEvStack          = 3
	traceEvGomaxprocs     = 4
	traceEvProcStart      = 5
	traceEvProcStop       = 6
	traceEvGCStart        = 7
	traceEvGCDone         = 8
	traceEvGCScanStart    = 9
	traceEvGCScanDone     = 10
	traceEvGCSweepStart   = 11
	traceEvGCSweepDone    = 12
	traceEvGoCreate       = 13
	traceEvGoStart        = 14
	traceEvGoEnd          = 15
	traceEvGoStop         = 16
	traceEvGoSched        = 17
	traceEvGoPreempt      = 18
	traceEvGoSleep        = 19
	traceEvGoBlock        = 20
	traceEvGoUnblock      = 21
	traceEvGoBlockSend    = 22
	traceEvGoBlockRecv    = 23
	traceEvGoBlockSelect  = 24
	traceEvGoBlockSync    = 25
	traceEvGoBlockCond    = 26
	traceEvGoBlockNet     = 27
	traceEvGoSysCall      = 28
	traceEvGoSysExit      = 29
	traceEvGoSysBlock     = 30
	traceEvGoWaiting      = 31
	traceEvGoInSyscall    = 32
	traceEvHeapAlloc      = 33
	traceEvNextGC         = 34
	traceEvTimerGoroutine = 35
	traceEvFutileWakeup   = 36
	traceEvCount          = 37
)

const (
	traceTickDiv = 16 + 48*(goarch_386|goarch_amd64|goarch_amd64p32)

	traceStackSize = 128

	traceGlobProc = -1

	traceBytesPerNumber = 10

	traceArgCountShift = 6

	traceFutileWakeup byte = 128
)

// traceBufHeader is per-P tracing buffer.
type traceBufHeader struct {
	link      *traceBuf
	lastSeq   uint64
	lastTicks uint64
	buf       []byte
	stk       [traceStackSize]uintptr
}

// traceBuf is per-P tracing buffer.
type traceBuf struct {
	traceBufHeader
	arr [64<<10 - unsafe.Sizeof(traceBufHeader{})]byte
}

// traceStackTable maps stack traces (arrays of PC's) to unique uint32 ids.
// It is lock-free for reading.
type traceStackTable struct {
	lock mutex
	seq  uint32
	mem  traceAlloc
	tab  [1 << 13]*traceStack
}

// traceStack is a single stack in traceStackTable.
type traceStack struct {
	link *traceStack
	hash uintptr
	id   uint32
	n    int
	stk  [0]uintptr
}

// traceAlloc is a non-thread-safe region allocator.
// It holds a linked list of traceAllocBlock.
type traceAlloc struct {
	head *traceAllocBlock
	off  uintptr
}

// traceAllocBlock is a block in traceAlloc.
type traceAllocBlock struct {
	next *traceAllocBlock
	data [64<<10 - ptrSize]byte
}
