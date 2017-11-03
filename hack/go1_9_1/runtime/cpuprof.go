// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// CPU profiling.
//
// The signal handler for the profiling clock tick adds a new stack trace
// to a log of recent traces. The log is read by a user goroutine that
// turns it into formatted profile data. If the reader does not keep up
// with the log, those writes will be recorded as a count of lost records.
// The actual profile buffer is in profbuf.go.

package runtime

const maxCPUProfStack = 64

type cpuProfile struct {
	lock mutex
	on   bool
	log  *profBuf

	extra     [1000]uintptr
	numExtra  int
	lostExtra uint64
}
