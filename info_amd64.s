// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

TEXT ·getg(SB), NOSPLIT, $0-8
    get_tls(CX)
    MOVQ    g(CX), AX
    MOVQ    AX, ret+0(FP)
    RET

TEXT ·get_goexit(SB), NOSPLIT, $0-8
    MOVQ    runtime·goexit(SB), AX
    MOVQ    AX, ret+0(FP)
    RET

TEXT ·find_goexit_pc(SB), NOSPLIT, $0-20
    MOVQ    offset+0(FP), DX
    ADDQ    runtime·goexit(SB), DX
    MOVQ    $0, R9

    get_tls(CX)
    MOVQ    g(CX), AX
    MOVQ    (g_stack+stack_hi)(AX), R8

    XORQ    CX, CX

    // Brute-force search goroutine stack from hi to SP.
    // Return the last match.
again:
    CMPQ    SP, R8
    JGE done
    SUBQ    $8, R8
    CMPQ    (R8), DX
    JNE again

    // Found runtime·goexit.
    MOVQ    R8, R9
    INCL    CX
    JMP again

done:
    MOVQ    R9, ret+8(FP)
    MOVL    CX, cnt+16(FP)
    RET

TEXT goroutine·goexit_jmp(SB), NOSPLIT, $0-0
    BYTE    $0x90	// NOP
    CALL	·hack_goexit(SB)	// does not return
    // traceback from goexit1 must hit code range of goexit
    BYTE	$0x90	// NOP

TEXT ·real_goexit(SB), NOSPLIT, $0-4
    MOVL    offset+0(FP), AX
    ADDQ    runtime·goexit(SB), AX
    CALL    AX
    // doesn't return
    BYTE	$0x90	// NOP

TEXT goroutine·goexit1(SB), NOSPLIT, $0-0
    CALL	runtime·goexit1(SB)

TEXT ·hook_goexit(SB), NOSPLIT, $0-5
    XORQ    BX, BX
    MOVL    offset+0(FP), BX

    // Calculate runtime·goexit offset.
    MOVQ    runtime·goexit(SB), DX
    ADDQ    BX, DX

    // Prepare ret position.
    ADDQ    goroutine·goexit_jmp(SB), BX

    get_tls(CX)
    MOVQ    g(CX), AX
    MOVQ    (g_stack+stack_hi)(AX), R8

    // Brute-force search goroutine stack from hi to SP.
    // Hack all entries.
again:
    SUBQ    $4, R8
    CMPQ    SP, R8
    JGE done
    CMPQ    (R8), DX
    JNE again

    // Found runtime·goexit in goroutine stack.
    // Replace it with goroutine·goexit_jmp.
    MOVQ    BX, (R8)

    // test
    SUBQ    $8, R8
    MOVQ    $0, (R8)
    ADDQ    $16, R8
    MOVQ    $0, (R8)

    MOVB    $1, ret+4(FP)
    RET

done:
    MOVB    $0, ret+4(FP)
    RET
