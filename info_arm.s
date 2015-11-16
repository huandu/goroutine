// Copyright 2015 Huan Du. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

#include "go_asm.h"
#include "textflag.h"

TEXT ·getg(SB), NOSPLIT, $0-4
    MOVW g, ret+0(FP)
    RET
