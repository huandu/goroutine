// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sys

const PtrSize = 4 << (^uintptr(0) >> 63)
const RegSize = 4 << (^Uintreg(0) >> 63)
const SpAlign = 1*(1-GoarchArm64) + 16*GoarchArm64
