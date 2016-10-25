package runtime

const (
	_SIGQUIT = 3
	_SIGSEGV = 11
	_SIGPROF = 27
)

type timespec struct {
	tv_sec  int64
	tv_nsec int32
}

type excregsarm struct {
	r0   uint32
	r1   uint32
	r2   uint32
	r3   uint32
	r4   uint32
	r5   uint32
	r6   uint32
	r7   uint32
	r8   uint32
	r9   uint32
	r10  uint32
	r11  uint32
	r12  uint32
	sp   uint32
	lr   uint32
	pc   uint32
	cpsr uint32
}

type exccontext struct {
	size                    uint32
	portable_context_offset uint32
	portable_context_size   uint32
	arch                    uint32
	regs_size               uint32
	reserved                [11]uint32
	regs                    excregsarm
}

type excportablecontext struct {
	pc uint32
	sp uint32
	fp uint32
}
