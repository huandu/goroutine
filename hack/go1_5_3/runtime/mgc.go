// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO(rsc): The code having to do with the heap bitmap needs very serious cleanup.
// It has gotten completely out of control.

// Garbage collector (GC).
//
// The GC runs concurrently with mutator threads, is type accurate (aka precise), allows multiple
// GC thread to run in parallel. It is a concurrent mark and sweep that uses a write barrier. It is
// non-generational and non-compacting. Allocation is done using size segregated per P allocation
// areas to minimize fragmentation while eliminating locks in the common case.
//
// The algorithm decomposes into several steps.
// This is a high level description of the algorithm being used. For an overview of GC a good
// place to start is Richard Jones' gchandbook.org.
//
// The algorithm's intellectual heritage includes Dijkstra's on-the-fly algorithm, see
// Edsger W. Dijkstra, Leslie Lamport, A. J. Martin, C. S. Scholten, and E. F. M. Steffens. 1978.
// On-the-fly garbage collection: an exercise in cooperation. Commun. ACM 21, 11 (November 1978),
// 966-975.
// For journal quality proofs that these steps are complete, correct, and terminate see
// Hudson, R., and Moss, J.E.B. Copying Garbage Collection without stopping the world.
// Concurrency and Computation: Practice and Experience 15(3-5), 2003.
//
//  0. Set phase = GCscan from GCoff.
//  1. Wait for all P's to acknowledge phase change.
//         At this point all goroutines have passed through a GC safepoint and
//         know we are in the GCscan phase.
//  2. GC scans all goroutine stacks, mark and enqueues all encountered pointers
//       (marking avoids most duplicate enqueuing but races may produce benign duplication).
//       Preempted goroutines are scanned before P schedules next goroutine.
//  3. Set phase = GCmark.
//  4. Wait for all P's to acknowledge phase change.
//  5. Now write barrier marks and enqueues black, grey, or white to white pointers.
//       Malloc still allocates white (non-marked) objects.
//  6. Meanwhile GC transitively walks the heap marking reachable objects.
//  7. When GC finishes marking heap, it preempts P's one-by-one and
//       retakes partial wbufs (filled by write barrier or during a stack scan of the goroutine
//       currently scheduled on the P).
//  8. Once the GC has exhausted all available marking work it sets phase = marktermination.
//  9. Wait for all P's to acknowledge phase change.
// 10. Malloc now allocates black objects, so number of unmarked reachable objects
//        monotonically decreases.
// 11. GC preempts P's one-by-one taking partial wbufs and marks all unmarked yet
//        reachable objects.
// 12. When GC completes a full cycle over P's and discovers no new grey
//         objects, (which means all reachable objects are marked) set phase = GCoff.
// 13. Wait for all P's to acknowledge phase change.
// 14. Now malloc allocates white (but sweeps spans before use).
//         Write barrier becomes nop.
// 15. GC does background sweeping, see description below.
// 16. When sufficient allocation has taken place replay the sequence starting at 0 above,
//         see discussion of GC rate below.

// Changing phases.
// Phases are changed by setting the gcphase to the next phase and possibly calling ackgcphase.
// All phase action must be benign in the presence of a change.
// Starting with GCoff
// GCoff to GCscan
//     GSscan scans stacks and globals greying them and never marks an object black.
//     Once all the P's are aware of the new phase they will scan gs on preemption.
//     This means that the scanning of preempted gs can't start until all the Ps
//     have acknowledged.
//     When a stack is scanned, this phase also installs stack barriers to
//     track how much of the stack has been active.
//     This transition enables write barriers because stack barriers
//     assume that writes to higher frames will be tracked by write
//     barriers. Technically this only needs write barriers for writes
//     to stack slots, but we enable write barriers in general.
// GCscan to GCmark
//     In GCmark, work buffers are drained until there are no more
//     pointers to scan.
//     No scanning of objects (making them black) can happen until all
//     Ps have enabled the write barrier, but that already happened in
//     the transition to GCscan.
// GCmark to GCmarktermination
//     The only change here is that we start allocating black so the Ps must acknowledge
//     the change before we begin the termination algorithm
// GCmarktermination to GSsweep
//     Object currently on the freelist must be marked black for this to work.
//     Are things on the free lists black or white? How does the sweep phase work?

// Concurrent sweep.
//
// The sweep phase proceeds concurrently with normal program execution.
// The heap is swept span-by-span both lazily (when a goroutine needs another span)
// and concurrently in a background goroutine (this helps programs that are not CPU bound).
// At the end of STW mark termination all spans are marked as "needs sweeping".
//
// The background sweeper goroutine simply sweeps spans one-by-one.
//
// To avoid requesting more OS memory while there are unswept spans, when a
// goroutine needs another span, it first attempts to reclaim that much memory
// by sweeping. When a goroutine needs to allocate a new small-object span, it
// sweeps small-object spans for the same object size until it frees at least
// one object. When a goroutine needs to allocate large-object span from heap,
// it sweeps spans until it frees at least that many pages into heap. There is
// one case where this may not suffice: if a goroutine sweeps and frees two
// nonadjacent one-page spans to the heap, it will allocate a new two-page
// span, but there can still be other one-page unswept spans which could be
// combined into a two-page span.
//
// It's critical to ensure that no operations proceed on unswept spans (that would corrupt
// mark bits in GC bitmap). During GC all mcaches are flushed into the central cache,
// so they are empty. When a goroutine grabs a new span into mcache, it sweeps it.
// When a goroutine explicitly frees an object or sets a finalizer, it ensures that
// the span is swept (either by sweeping it, or by waiting for the concurrent sweep to finish).
// The finalizer goroutine is kicked off only when all spans are swept.
// When the next GC starts, it sweeps all not-yet-swept spans (if any).

// GC rate.
// Next GC is after we've allocated an extra amount of memory proportional to
// the amount already in use. The proportion is controlled by GOGC environment variable
// (100 by default). If GOGC=100 and we're using 4M, we'll GC again when we get to 8M
// (this mark is tracked in next_gc variable). This keeps the GC cost in linear
// proportion to the allocation cost. Adjusting GOGC just changes the linear constant
// (and also the amount of extra memory used).

package runtime

const (
	_DebugGC         = 0
	_ConcurrentSweep = true
	_FinBlockSize    = 4 * 1024
	_RootData        = 0
	_RootBss         = 1
	_RootFinalizers  = 2
	_RootSpans       = 3
	_RootFlushCaches = 4
	_RootCount       = 5

	debugStackBarrier = false

	sweepMinHeapDistance = 1024 * 1024
)

// defaultHeapMinimum is the value of heapminimum for GOGC==100.
const defaultHeapMinimum = 4 << 20

const (
	_GCoff = iota
	_GCstw
	_GCscan
	_GCmark
	_GCmarktermination
)

// gcMarkWorkerMode represents the mode that a concurrent mark worker
// should operate in.
//
// Concurrent marking happens through four different mechanisms. One
// is mutator assists, which happen in response to allocations and are
// not scheduled. The other three are variations in the per-P mark
// workers and are distinguished by gcMarkWorkerMode.
type gcMarkWorkerMode int

const (
	gcMarkWorkerDedicatedMode gcMarkWorkerMode = iota

	gcMarkWorkerFractionalMode

	gcMarkWorkerIdleMode
)

type gcControllerState struct {
	scanWork int64

	bgScanCredit int64

	assistTime int64

	dedicatedMarkTime int64

	fractionalMarkTime int64

	idleMarkTime int64

	bgMarkStartTime int64

	assistStartTime int64

	heapGoal uint64

	dedicatedMarkWorkersNeeded int64

	assistRatio float64

	fractionalUtilizationGoal float64

	triggerRatio float64

	_ [_CacheLineSize]byte

	fractionalMarkWorkersNeeded int64

	_ [_CacheLineSize]byte
}

// gcGoalUtilization is the goal CPU utilization for background
// marking as a fraction of GOMAXPROCS.
const gcGoalUtilization = 0.25

// gcBgCreditSlack is the amount of scan work credit background
// scanning can accumulate locally before updating
// gcController.bgScanCredit. Lower values give mutator assists more
// accurate accounting of background scanning. Higher values reduce
// memory contention.
const gcBgCreditSlack = 2000

// gcAssistTimeSlack is the nanoseconds of mutator assist time that
// can accumulate on a P before updating gcController.assistTime.
const gcAssistTimeSlack = 5000

// bgMarkSignal synchronizes the GC coordinator and background mark workers.
type bgMarkSignal struct {
	done uint32

	lock mutex
	g    *g
	wake bool
}

const (
	gcBackgroundMode = iota
	gcForceMode
	gcForceBlockMode
)
