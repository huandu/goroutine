// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

const (
	deferHeaderSize = unsafe.Sizeof(_defer{})
	minDeferAlloc   = (deferHeaderSize + 15) &^ 15
	minDeferArgs    = minDeferAlloc - deferHeaderSize
)
