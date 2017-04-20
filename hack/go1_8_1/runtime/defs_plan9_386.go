package runtime

const _PAGESIZE = 0x1000

type ureg struct {
	di    uint32
	si    uint32
	bp    uint32
	nsp   uint32
	bx    uint32
	dx    uint32
	cx    uint32
	ax    uint32
	gs    uint32
	fs    uint32
	es    uint32
	ds    uint32
	trap  uint32
	ecode uint32
	pc    uint32
	cs    uint32
	flags uint32
	sp    uint32
	ss    uint32
}

type sigctxt struct {
	u *ureg
}
