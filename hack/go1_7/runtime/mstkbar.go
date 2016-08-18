// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Garbage collector: stack barriers
//
// Stack barriers enable the garbage collector to determine how much
// of a gorountine stack has changed between when a stack is scanned
// during the concurrent scan phase and when it is re-scanned during
// the stop-the-world mark termination phase. Mark termination only
// needs to re-scan the changed part, so for deep stacks this can
// significantly reduce GC pause time compared to the alternative of
// re-scanning whole stacks. The deeper the stacks, the more stack
// barriers help.
//
// When stacks are scanned during the concurrent scan phase, the stack
// scan installs stack barriers by selecting stack frames and
// overwriting the saved return PCs (or link registers) of these
// frames with the PC of a "stack barrier trampoline". Later, when a
// selected frame returns, it "returns" to this trampoline instead of
// returning to its actual caller. The trampoline records that the
// stack has unwound past this frame and jumps to the original return
// PC recorded when the stack barrier was installed. Mark termination
// re-scans only as far as the first frame that hasn't hit a stack
// barrier and then removes and un-hit stack barriers.
//
// This scheme is very lightweight. No special code is required in the
// mutator to record stack unwinding and the trampoline is only a few
// assembly instructions.
//
// Book-keeping
// ------------
//
// The primary cost of stack barriers is book-keeping: the runtime has
// to record the locations of all stack barriers and the original
// return PCs in order to return to the correct caller when a stack
// barrier is hit and so it can remove un-hit stack barriers. In order
// to minimize this cost, the Go runtime places stack barriers in
// exponentially-spaced frames, starting 1K past the current frame.
// The book-keeping structure hence grows logarithmically with the
// size of the stack and mark termination re-scans at most twice as
// much stack as necessary.
//
// The runtime reserves space for this book-keeping structure at the
// top of the stack allocation itself (just above the outermost
// frame). This is necessary because the regular memory allocator can
// itself grow the stack, and hence can't be used when allocating
// stack-related structures.
//
// For debugging, the runtime also supports installing stack barriers
// at every frame. However, this requires significantly more
// book-keeping space.
//
// Correctness
// -----------
//
// The runtime and the compiler cooperate to ensure that all objects
// reachable from the stack as of mark termination are marked.
// Anything unchanged since the concurrent scan phase will be marked
// because it is marked by the concurrent scan. After the concurrent
// scan, there are three possible classes of stack modifications that
// must be tracked:
//
// 1) Mutator writes below the lowest un-hit stack barrier. This
// includes all writes performed by an executing function to its own
// stack frame. This part of the stack will be re-scanned by mark
// termination, which will mark any objects made reachable from
// modifications to this part of the stack.
//
// 2) Mutator writes above the lowest un-hit stack barrier. It's
// possible for a mutator to modify the stack above the lowest un-hit
// stack barrier if a higher frame has passed down a pointer to a
// stack variable in its frame. This is called an "up-pointer". The
// compiler ensures that writes through up-pointers have an
// accompanying write barrier (it simply doesn't distinguish between
// writes through up-pointers and writes through heap pointers). This
// write barrier marks any object made reachable from modifications to
// this part of the stack.
//
// 3) Runtime writes to the stack. Various runtime operations such as
// sends to unbuffered channels can write to arbitrary parts of the
// stack, including above the lowest un-hit stack barrier. We solve
// this in two ways. In many cases, the runtime can perform an
// explicit write barrier operation like in case 2. However, in the
// case of bulk memory move (typedmemmove), the runtime doesn't
// necessary have ready access to a pointer bitmap for the memory
// being copied, so it simply unwinds any stack barriers below the
// destination.
//
// Gotchas
// -------
//
// Anything that inspects or manipulates the stack potentially needs
// to understand stack barriers. The most obvious case is that
// gentraceback needs to use the original return PC when it encounters
// the stack barrier trampoline. Anything that unwinds the stack such
// as panic/recover must unwind stack barriers in tandem with
// unwinding the stack.
//
// Stack barriers require that any goroutine whose stack has been
// scanned must execute write barriers. Go solves this by simply
// enabling write barriers globally during the concurrent scan phase.
// However, traditionally, write barriers are not enabled during this
// phase.
//
// Synchronization
// ---------------
//
// For the most part, accessing and modifying stack barriers is
// synchronized around GC safe points. Installing stack barriers
// forces the G to a safe point, while all other operations that
// modify stack barriers run on the G and prevent it from reaching a
// safe point.
//
// Subtlety arises when a G may be tracebacked when *not* at a safe
// point. This happens during sigprof. For this, each G has a "stack
// barrier lock" (see gcLockStackBarriers, gcUnlockStackBarriers).
// Operations that manipulate stack barriers acquire this lock, while
// sigprof tries to acquire it and simply skips the traceback if it
// can't acquire it. There is one exception for performance and
// complexity reasons: hitting a stack barrier manipulates the stack
// barrier list without acquiring the stack barrier lock. For this,
// gentraceback performs a special fix up if the traceback starts in
// the stack barrier function.

package runtime

const debugStackBarrier = false
