// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

const (
	_AT_RANDOM       = 25
	_AT_SYSINFO_EHDR = 33
	_AT_NULL         = 0

	_PT_LOAD    = 1
	_PT_DYNAMIC = 2

	_DT_NULL   = 0
	_DT_HASH   = 4
	_DT_STRTAB = 5
	_DT_SYMTAB = 6
	_DT_VERSYM = 0x6ffffff0
	_DT_VERDEF = 0x6ffffffc

	_VER_FLG_BASE = 0x1

	_SHN_UNDEF = 0

	_SHT_DYNSYM = 11

	_STT_FUNC = 2

	_STB_GLOBAL = 1
	_STB_WEAK   = 2

	_EI_NIDENT = 16
)

type elf64Sym struct {
	st_name  uint32
	st_info  byte
	st_other byte
	st_shndx uint16
	st_value uint64
	st_size  uint64
}

type elf64Verdef struct {
	vd_version uint16
	vd_flags   uint16
	vd_ndx     uint16
	vd_cnt     uint16
	vd_hash    uint32
	vd_aux     uint32
	vd_next    uint32
}

type elf64Ehdr struct {
	e_ident     [_EI_NIDENT]byte
	e_type      uint16
	e_machine   uint16
	e_version   uint32
	e_entry     uint64
	e_phoff     uint64
	e_shoff     uint64
	e_flags     uint32
	e_ehsize    uint16
	e_phentsize uint16
	e_phnum     uint16
	e_shentsize uint16
	e_shnum     uint16
	e_shstrndx  uint16
}

type elf64Phdr struct {
	p_type   uint32
	p_flags  uint32
	p_offset uint64
	p_vaddr  uint64
	p_paddr  uint64
	p_filesz uint64
	p_memsz  uint64
	p_align  uint64
}

type elf64Shdr struct {
	sh_name      uint32
	sh_type      uint32
	sh_flags     uint64
	sh_addr      uint64
	sh_offset    uint64
	sh_size      uint64
	sh_link      uint32
	sh_info      uint32
	sh_addralign uint64
	sh_entsize   uint64
}

type elf64Dyn struct {
	d_tag int64
	d_val uint64
}

type elf64Verdaux struct {
	vda_name uint32
	vda_next uint32
}

type elf64Auxv struct {
	a_type uint64
	a_val  uint64
}

type symbol_key struct {
	name     string
	sym_hash uint32
	ptr      *uintptr
}

type version_key struct {
	version  string
	ver_hash uint32
}

type vdso_info struct {
	valid bool

	load_addr   uintptr
	load_offset uintptr

	symtab     *[1 << 32]elf64Sym
	symstrings *[1 << 32]byte
	chain      []uint32
	bucket     []uint32

	versym *[1 << 32]uint16
	verdef *elf64Verdef
}
