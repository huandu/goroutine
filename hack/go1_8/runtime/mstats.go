// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Memory statistics

package runtime

// Statistics.
// If you edit this structure, also edit type MemStats below.
// Their layouts must match exactly.
//
// For detailed descriptions see the documentation for MemStats.
// Fields that differ from MemStats are further documented here.
//
// Many of these fields are updated on the fly, while others are only
// updated when updatememstats is called.
type mstats struct {
	alloc       uint64
	total_alloc uint64
	sys         uint64
	nlookup     uint64
	nmalloc     uint64
	nfree       uint64

	heap_alloc    uint64
	heap_sys      uint64
	heap_idle     uint64
	heap_inuse    uint64
	heap_released uint64
	heap_objects  uint64

	stacks_inuse uint64
	stacks_sys   uint64
	mspan_inuse  uint64
	mspan_sys    uint64
	mcache_inuse uint64
	mcache_sys   uint64
	buckhash_sys uint64
	gc_sys       uint64
	other_sys    uint64

	next_gc         uint64
	last_gc         uint64
	pause_total_ns  uint64
	pause_ns        [256]uint64
	pause_end       [256]uint64
	numgc           uint32
	numforcedgc     uint32
	gc_cpu_fraction float64
	enablegc        bool
	debuggc         bool

	by_size [_NumSizeClasses]struct {
		size    uint32
		nmalloc uint64
		nfree   uint64
	}

	tinyallocs uint64

	gc_trigger uint64

	heap_live uint64

	heap_scan uint64

	heap_marked uint64
}

// A MemStats records statistics about the memory allocator.
type MemStats struct {
	Alloc uint64

	TotalAlloc uint64

	Sys uint64

	Lookups uint64

	Mallocs uint64

	Frees uint64

	HeapAlloc uint64

	HeapSys uint64

	HeapIdle uint64

	HeapInuse uint64

	HeapReleased uint64

	HeapObjects uint64

	StackInuse uint64

	StackSys uint64

	MSpanInuse uint64

	MSpanSys uint64

	MCacheInuse uint64

	MCacheSys uint64

	BuckHashSys uint64

	GCSys uint64

	OtherSys uint64

	NextGC uint64

	LastGC uint64

	PauseTotalNs uint64

	PauseNs [256]uint64

	PauseEnd [256]uint64

	NumGC uint32

	NumForcedGC uint32

	GCCPUFraction float64

	EnableGC bool

	DebugGC bool

	BySize [61]struct {
		Size uint32

		Mallocs uint64

		Frees uint64
	}
}
