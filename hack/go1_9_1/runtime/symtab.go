// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// Frames may be used to get function/file/line information for a
// slice of PC values returned by Callers.
type Frames struct {
	callers []uintptr

	stackExpander stackExpander
}

// Frame is the information returned by Frames for each call frame.
type Frame struct {
	PC uintptr

	Func *Func

	Function string

	File string
	Line int

	Entry uintptr
}

// stackExpander expands a call stack of PCs into a sequence of
// Frames. It tracks state across PCs necessary to perform this
// expansion.
//
// This is the core of the Frames implementation, but is a separate
// internal API to make it possible to use within the runtime without
// heap-allocating the PC slice. The only difference with the public
// Frames API is that the caller is responsible for threading the PC
// slice between expansion steps in this API. If escape analysis were
// smarter, we may not need this (though it may have to be a lot
// smarter).
type stackExpander struct {
	pcExpander pcExpander

	wasPanic bool

	skip int
}

// A pcExpander expands a single PC into a sequence of Frames.
type pcExpander struct {
	more bool

	pc uintptr

	frames []Frame

	funcInfo funcInfo

	inlTree *[1 << 20]inlinedCall

	file string
	line int32

	inlIndex int32
}

// A Func represents a Go function in the running binary.
type Func struct {
	opaque struct{}
}

// PCDATA and FUNCDATA table indexes.
//
// See funcdata.h and ../cmd/internal/obj/funcdata.go.
const (
	_PCDATA_StackMapIndex       = 0
	_PCDATA_InlTreeIndex        = 1
	_FUNCDATA_ArgsPointerMaps   = 0
	_FUNCDATA_LocalsPointerMaps = 1
	_FUNCDATA_InlTree           = 2
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
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*_type

	next *moduledata
}

// A modulehash is used to compare the ABI of a new module or a
// package in a new module with the loaded program.
//
// For each shared library a module links against, the linker creates an entry in the
// moduledata.modulehashes slice containing the name of the module, the abi hash seen
// at link time and a pointer to the runtime abi hash. These are checked in
// moduledataverify1 below.
//
// For each loaded plugin, the pkghashes slice has a modulehash of the
// newly loaded package that can be used to check the plugin's version of
// a package against any previously loaded version of the package.
// This is done in plugin.lastmoduleinit.
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

type functab struct {
	entry   uintptr
	funcoff uintptr
}

type textsect struct {
	vaddr    uintptr
	length   uintptr
	baseaddr uintptr
}

const minfunc = 16
const pcbucketsize = 256 * minfunc

// findfunctab is an array of these structures.
// Each bucket represents 4096 bytes of the text segment.
// Each subbucket represents 256 bytes of the text segment.
// To find a function given a pc, locate the bucket and subbucket for
// that pc. Add together the idx and subbucket value to obtain a
// function index. Then scan the functab array starting at that
// index to find the target function.
// This table uses 20 bytes for every 4096 bytes of code, or ~0.5% overhead.
type findfuncbucket struct {
	idx        uint32
	subbuckets [16]byte
}

const debugPcln = false

type funcInfo struct {
	*_func
	datap *moduledata
}

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

// inlinedCall is the encoding of entries in the FUNCDATA_InlTree table.
type inlinedCall struct {
	parent int32
	file   int32
	line   int32
	func_  int32
}
