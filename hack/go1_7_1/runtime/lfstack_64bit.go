// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 arm64 mips64 mips64le ppc64 ppc64le s390x

package runtime

const (
	addrBits = 48

	cntBits = 64 - addrBits + 3
)
