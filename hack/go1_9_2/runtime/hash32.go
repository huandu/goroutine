// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hashing algorithm inspired by
//   xxhash: https://code.google.com/p/xxhash/
// cityhash: https://code.google.com/p/cityhash/

// +build 386 arm mips mipsle

package runtime

const (
	m1 = 3168982561
	m2 = 3339683297
	m3 = 832293441
	m4 = 2336365089
)
