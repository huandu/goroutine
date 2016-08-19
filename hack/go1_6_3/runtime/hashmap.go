// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_6_3/runtime/internal/sys"
	"unsafe"
)

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits

	loadFactor = 6.5

	maxKeySize   = 128
	maxValueSize = 128

	dataOffset = unsafe.Offsetof(struct {
		b bmap
		v int64
	}{}.v)

	empty          = 0
	evacuatedEmpty = 1
	evacuatedX     = 2
	evacuatedY     = 3
	minTopHash     = 4

	iterator    = 1
	oldIterator = 2
	hashWriting = 4

	noCheck = 1<<(8*sys.PtrSize) - 1
)

// A header for a Go map.
type hmap struct {
	count int
	flags uint8
	B     uint8
	hash0 uint32

	buckets    unsafe.Pointer
	oldbuckets unsafe.Pointer
	nevacuate  uintptr

	overflow *[2]*[]*bmap
}

// A bucket for a Go map.
type bmap struct {
	tophash [bucketCnt]uint8
}

// A hash iteration structure.
// If you modify hiter, also change cmd/internal/gc/reflect.go to indicate
// the layout of this structure.
type hiter struct {
	key         unsafe.Pointer
	value       unsafe.Pointer
	t           *maptype
	h           *hmap
	buckets     unsafe.Pointer
	bptr        *bmap
	overflow    [2]*[]*bmap
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

const initialZeroSize = 1024
