// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const ptrSize = 4 << (^uintptr(0) >> 63)
const regSize = 4 << (^uintreg(0) >> 63)
const spAlign = 1*(1-goarch_arm64) + 16*goarch_arm64

type neverCallThisFunction struct{}

// argp used in Defer structs when there is no argp.
const _NoArgs = ^uintptr(0)
