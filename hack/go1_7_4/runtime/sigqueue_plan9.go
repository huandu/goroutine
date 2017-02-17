// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements runtime support for signal handling.

package runtime

const qsize = 64

type noteData struct {
	s [_ERRMAX]byte
	n int
}

type noteQueue struct {
	lock mutex
	data [qsize]noteData
	ri   int
	wi   int
	full bool
}
