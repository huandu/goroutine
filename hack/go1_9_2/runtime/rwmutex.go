// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// An rwmutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// This is a variant of sync.RWMutex, for the runtime package.
// Like mutex, rwmutex blocks the calling M.
// It does not interact with the goroutine scheduler.
type rwmutex struct {
	rLock      mutex
	readers    muintptr
	readerPass uint32

	wLock  mutex
	writer muintptr

	readerCount uint32
	readerWait  uint32
}

const rwmutexMaxReaders = 1 << 30
