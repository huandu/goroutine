// Copyright 2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

// A profBuf is a lock-free buffer for profiling events,
// safe for concurrent use by one reader and one writer.
// The writer may be a signal handler running without a user g.
// The reader is assumed to be a user g.
//
// Each logged event corresponds to a fixed size header, a list of
// uintptrs (typically a stack), and exactly one unsafe.Pointer tag.
// The header and uintptrs are stored in the circular buffer data and the
// tag is stored in a circular buffer tags, running in parallel.
// In the circular buffer data, each event takes 2+hdrsize+len(stk)
// words: the value 2+hdrsize+len(stk), then the time of the event, then
// hdrsize words giving the fixed-size header, and then len(stk) words
// for the stack.
//
// The current effective offsets into the tags and data circular buffers
// for reading and writing are stored in the high 30 and low 32 bits of r and w.
// The bottom bits of the high 32 are additional flag bits in w, unused in r.
// "Effective" offsets means the total number of reads or writes, mod 2^length.
// The offset in the buffer is the effective offset mod the length of the buffer.
// To make wraparound mod 2^length match wraparound mod length of the buffer,
// the length of the buffer must be a power of two.
//
// If the reader catches up to the writer, a flag passed to read controls
// whether the read blocks until more data is available. A read returns a
// pointer to the buffer data itself; the caller is assumed to be done with
// that data at the next read. The read offset rNext tracks the next offset to
// be returned by read. By definition, r ≤ rNext ≤ w (before wraparound),
// and rNext is only used by the reader, so it can be accessed without atomics.
//
// If the writer gets ahead of the reader, so that the buffer fills,
// future writes are discarded and replaced in the output stream by an
// overflow entry, which has size 2+hdrsize+1, time set to the time of
// the first discarded write, a header of all zeroed words, and a "stack"
// containing one word, the number of discarded writes.
//
// Between the time the buffer fills and the buffer becomes empty enough
// to hold more data, the overflow entry is stored as a pending overflow
// entry in the fields overflow and overflowTime. The pending overflow
// entry can be turned into a real record by either the writer or the
// reader. If the writer is called to write a new record and finds that
// the output buffer has room for both the pending overflow entry and the
// new record, the writer emits the pending overflow entry and the new
// record into the buffer. If the reader is called to read data and finds
// that the output buffer is empty but that there is a pending overflow
// entry, the reader will return a synthesized record for the pending
// overflow entry.
//
// Only the writer can create or add to a pending overflow entry, but
// either the reader or the writer can clear the pending overflow entry.
// A pending overflow entry is indicated by the low 32 bits of 'overflow'
// holding the number of discarded writes, and overflowTime holding the
// time of the first discarded write. The high 32 bits of 'overflow'
// increment each time the low 32 bits transition from zero to non-zero
// or vice versa. This sequence number avoids ABA problems in the use of
// compare-and-swap to coordinate between reader and writer.
// The overflowTime is only written when the low 32 bits of overflow are
// zero, that is, only when there is no pending overflow entry, in
// preparation for creating a new one. The reader can therefore fetch and
// clear the entry atomically using
//
//	for {
//		overflow = load(&b.overflow)
//		if uint32(overflow) == 0 {
//			// no pending entry
//			break
//		}
//		time = load(&b.overflowTime)
//		if cas(&b.overflow, overflow, ((overflow>>32)+1)<<32) {
//			// pending entry cleared
//			break
//		}
//	}
//	if uint32(overflow) > 0 {
//		emit entry for uint32(overflow), time
//	}
//
type profBuf struct {
	r, w         profAtomic
	overflow     uint64
	overflowTime uint64
	eof          uint32

	hdrsize uintptr
	data    []uint64
	tags    []unsafe.Pointer

	rNext       profIndex
	overflowBuf []uint64
	wait        note
}

// A profAtomic is the atomically-accessed word holding a profIndex.
type profAtomic uint64

// A profIndex is the packet tag and data counts and flags bits, described above.
type profIndex uint64

const (
	profReaderSleeping profIndex = 1 << 32
	profWriteExtra     profIndex = 1 << 33
)

// profBufReadMode specifies whether to block when no data is available to read.
type profBufReadMode int

const (
	profBufBlocking profBufReadMode = iota
	profBufNonBlocking
)
