// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code to check that pointer writes follow the cgo rules.
// These functions are invoked via the write barrier when debug.cgocheck > 1.

package runtime

const cgoWriteBarrierFail = "Go pointer stored into non-Go memory"
