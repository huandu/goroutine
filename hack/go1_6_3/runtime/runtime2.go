// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/huandu/goroutine/hack/go1_6_3/runtime/internal/sys"
	"unsafe"
)

/*
 * defined constants
 */
const (
	_Gidle = iota
	_Grunnable
	_Grunning
	_Gsyscall
	_Gwaiting
	_Gmoribund_unused
	_Gdead
	_Genqueue
	_Gcopystack

	_Gscan = 0x1000

	_Gscanrunnable = _Gscan + _Grunnable
	_Gscanrunning  = _Gscan + _Grunning
	_Gscansyscall  = _Gscan + _Gsyscall
	_Gscanwaiting  = _Gscan + _Gwaiting

	_Gscanenqueue = _Gscan + _Genqueue
)

const (
	_Pidle = iota
	_Prunning
	_Psyscall
	_Pgcstop
	_Pdead
)

type mutex struct {
	key uintptr
}

type note struct {
	key uintptr
}

type funcval struct {
	fn uintptr
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

// A guintptr holds a goroutine pointer, but typed as a uintptr
// to bypass write barriers. It is used in the Gobuf goroutine state
// and in scheduling lists that are manipulated without a P.
//
// The Gobuf.g goroutine pointer is almost always updated by assembly code.
// In one of the few places it is updated by Go code - func save - it must be
// treated as a uintptr to avoid a write barrier being emitted at a bad time.
// Instead of figuring out how to emit the write barriers missing in the
// assembly manipulation, we change the type of the field to uintptr,
// so that it does not require write barriers at all.
//
// Goroutine structs are published in the allg list and never freed.
// That will keep the goroutine structs from being collected.
// There is never a time that Gobuf.g's contain the only references
// to a goroutine: the publishing of the goroutine in allg comes first.
// Goroutine pointers are also kept in non-GC-visible places like TLS,
// so I can't see them ever moving. If we did want to start moving data
// in the GC, we'd need to allocate the goroutine structs from an
// alternate arena. Using guintptr doesn't make that problem any worse.
type guintptr uintptr

type puintptr uintptr

type muintptr uintptr

type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer
	ret  sys.Uintreg
	lr   uintptr
	bp   uintptr
}

// Known to compiler.
// Changes here must also be made in src/cmd/internal/gc/select.go's selecttype.
type sudog struct {
	g           *g
	selectdone  *uint32
	next        *sudog
	prev        *sudog
	elem        unsafe.Pointer
	releasetime int64
	nrelease    int32
	waitlink    *sudog
}

type gcstats struct {
	nhandoff    uint64
	nhandoffcnt uint64
	nprocyield  uint64
	nosyield    uint64
	nsleep      uint64
}

type libcall struct {
	fn   uintptr
	n    uintptr
	args uintptr
	r1   uintptr
	r2   uintptr
	err  uintptr
}

// describes how to handle callback
type wincallbackcontext struct {
	gobody       unsafe.Pointer
	argsize      uintptr
	restorestack uintptr
	cleanstack   bool
}

// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}

// stkbar records the state of a G's stack barrier.
type stkbar struct {
	savedLRPtr uintptr
	savedLRVal uintptr
}

type g struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr

	_panic         *_panic
	_defer         *_defer
	m              *m
	stackAlloc     uintptr
	sched          gobuf
	syscallsp      uintptr
	syscallpc      uintptr
	stkbar         []stkbar
	stkbarPos      uintptr
	stktopsp       uintptr
	param          unsafe.Pointer
	atomicstatus   uint32
	stackLock      uint32
	goid           int64
	waitsince      int64
	waitreason     string
	schedlink      guintptr
	preempt        bool
	paniconfault   bool
	preemptscan    bool
	gcscandone     bool
	gcscanvalid    bool
	throwsplit     bool
	raceignore     int8
	sysblocktraced bool
	sysexitticks   int64
	sysexitseq     uint64
	lockedm        *m
	sig            uint32
	writebuf       []byte
	sigcode0       uintptr
	sigcode1       uintptr
	sigpc          uintptr
	gopc           uintptr
	startpc        uintptr
	racectx        uintptr
	waiting        *sudog

	gcAssistBytes int64
}

type m struct {
	g0      *g
	morebuf gobuf
	divmod  uint32

	procid        uint64
	gsignal       *g
	sigmask       sigset
	tls           [6]uintptr
	mstartfn      func()
	curg          *g
	caughtsig     guintptr
	p             puintptr
	nextp         puintptr
	id            int32
	mallocing     int32
	throwing      int32
	preemptoff    string
	locks         int32
	softfloat     int32
	dying         int32
	profilehz     int32
	helpgc        int32
	spinning      bool
	blocked       bool
	inwb          bool
	newSigstack   bool
	printlock     int8
	fastrand      uint32
	ncgocall      uint64
	ncgo          int32
	park          note
	alllink       *m
	schedlink     muintptr
	machport      uint32
	mcache        *mcache
	lockedg       *g
	createstack   [32]uintptr
	freglo        [16]uint32
	freghi        [16]uint32
	fflag         uint32
	locked        uint32
	nextwaitm     uintptr
	gcstats       gcstats
	needextram    bool
	traceback     uint8
	waitunlockf   unsafe.Pointer
	waitlock      unsafe.Pointer
	waittraceev   byte
	waittraceskip int
	startingtrace bool
	syscalltick   uint32

	thread uintptr

	libcall   libcall
	libcallpc uintptr
	libcallsp uintptr
	libcallg  guintptr
	syscall   libcall

	mOS
}

type p struct {
	lock mutex

	id          int32
	status      uint32
	link        puintptr
	schedtick   uint32
	syscalltick uint32
	m           muintptr
	mcache      *mcache

	deferpool    [5][]*_defer
	deferpoolbuf [5][32]*_defer

	goidcache    uint64
	goidcacheend uint64

	runqhead uint32
	runqtail uint32
	runq     [256]guintptr

	runnext guintptr

	gfree    *g
	gfreecnt int32

	sudogcache []*sudog
	sudogbuf   [128]*sudog

	tracebuf traceBufPtr

	palloc persistentAlloc

	gcAssistTime     int64
	gcBgMarkWorker   guintptr
	gcMarkWorkerMode gcMarkWorkerMode

	gcw gcWork

	runSafePointFn uint32

	pad [64]byte
}

const (
	_MaxGomaxprocs = 1 << 8
)

type schedt struct {
	goidgen  uint64
	lastpoll uint64

	lock mutex

	midle        muintptr
	nmidle       int32
	nmidlelocked int32
	mcount       int32
	maxmcount    int32

	ngsys uint32

	pidle      puintptr
	npidle     uint32
	nmspinning uint32

	runqhead guintptr
	runqtail guintptr
	runqsize int32

	gflock mutex
	gfree  *g
	ngfree int32

	sudoglock  mutex
	sudogcache *sudog

	deferlock mutex
	deferpool [5]*_defer

	gcwaiting  uint32
	stopwait   int32
	stopnote   note
	sysmonwait uint32
	sysmonnote note

	safePointFn   func(*p)
	safePointWait int32
	safePointNote note

	profilehz int32

	procresizetime int64
	totaltime      int64
}

// The m->locked word holds two pieces of state counting active calls to LockOSThread/lockOSThread.
// The low bit (LockExternal) is a boolean reporting whether any LockOSThread call is active.
// External locks are not recursive; a second lock is silently ignored.
// The upper bits of m->locked record the nesting depth of calls to lockOSThread
// (counting up by LockInternal), popped by unlockOSThread (counting down by LockInternal).
// Internal locks can be recursive. For instance, a lock for cgo can occur while the main
// goroutine is holding the lock during the initialization phase.
const (
	_LockExternal = 1
	_LockInternal = 2
)

type sigtabtt struct {
	flags int32
	name  *int8
}

const (
	_SigNotify = 1 << iota
	_SigKill
	_SigThrow
	_SigPanic
	_SigDefault
	_SigHandling
	_SigGoExit
	_SigSetStack
	_SigUnblock
)

// Layout of in-memory per-function information prepared by linker
// See https://golang.org/s/go12symtab.
// Keep in sync with linker
// and with package debug/gosym and with symtab.go in package runtime.
type _func struct {
	entry   uintptr
	nameoff int32

	args int32
	_    int32

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	unused int32
	fun    [1]uintptr
}

// Lock-free stack node.
// // Also known to export_test.go.
type lfnode struct {
	next    uint64
	pushcnt uintptr
}

type forcegcstate struct {
	lock mutex
	g    *g
	idle uint32
}

/*
 * known to compiler
 */
const (
	_Structrnd = sys.RegSize
)

/*
 * deferred subroutine calls
 */
type _defer struct {
	siz     int32
	started bool
	sp      uintptr
	pc      uintptr
	fn      *funcval
	_panic  *_panic
	link    *_defer
}

/*
 * panics
 */
type _panic struct {
	argp      unsafe.Pointer
	arg       interface{}
	link      *_panic
	recovered bool
	aborted   bool
}

type stkframe struct {
	fn       *_func
	pc       uintptr
	continpc uintptr
	lr       uintptr
	sp       uintptr
	fp       uintptr
	varp     uintptr
	argp     uintptr
	arglen   uintptr
	argmap   *bitvector
}

const (
	_TraceRuntimeFrames = 1 << iota
	_TraceTrap
	_TraceJumpStack
)

const (
	_TracebackMaxFrames = 100
)
