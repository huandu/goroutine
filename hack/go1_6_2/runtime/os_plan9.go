// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

type mOS struct {
	waitsemacount uint32
	notesig       *int8
	errstr        *byte
}

type _Plink uintptr
