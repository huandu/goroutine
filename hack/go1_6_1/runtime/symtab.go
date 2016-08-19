// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// A Func represents a Go function in the running binary.
type Func struct {
	opaque struct{}
}

// funcdata.h
const (
	_PCDATA_StackMapIndex       = 0
	_FUNCDATA_ArgsPointerMaps   = 0
	_FUNCDATA_LocalsPointerMaps = 1
	_ArgsSizeUnknown            = -0x80000000
)

// moduledata records information about the layout of the executable
// image. It is written by the linker. Any changes here must be
// matched changes to the code in cmd/internal/ld/symtab.go:symtab.
// moduledata is stored in read-only memory; none of the pointers here
// are visible to the garbage collector.
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr

	typelinks []*_type

	modulename   string
	modulehashes []modulehash

	gcdatamask, gcbssmask bitvector

	next *moduledata
}

// For each shared library a module links against, the linker creates an entry in the
// moduledata.modulehashes slice containing the name of the module, the abi hash seen
// at link time and a pointer to the runtime abi hash. These are checked in
// moduledataverify1 below.
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

type functab struct {
	entry   uintptr
	funcoff uintptr
}

const minfunc = 16
const pcbucketsize = 256 * minfunc

// findfunctab is an array of these structures.
// Each bucket represents 4096 bytes of the text segment.
// Each subbucket represents 256 bytes of the text segment.
// To find a function given a pc, locate the bucket and subbucket for
// that pc.  Add together the idx and subbucket value to obtain a
// function index.  Then scan the functab array starting at that
// index to find the target function.
// This table uses 20 bytes for every 4096 bytes of code, or ~0.5% overhead.
type findfuncbucket struct {
	idx        uint32
	subbuckets [16]byte
}

const debugPcln = false

type pcvalueCache struct {
	entries [16]pcvalueCacheEnt
}

type pcvalueCacheEnt struct {
	targetpc uintptr
	off      int32

	val int32
}

type stackmap struct {
	n        int32
	nbit     int32
	bytedata [1]byte
}
