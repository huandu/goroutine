// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type mts struct {
	tv_sec  int64
	tv_nsec int64
}

type mscratch struct {
	v [6]uintptr
}

type mOS struct {
	waitsema uintptr
	perrno   *int32

	ts      mts
	scratch mscratch
}

type libcFunc uintptr
