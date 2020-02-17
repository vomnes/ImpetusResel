package net

import (
	"syscall"
)

const (
	FDSize = 1024
)

// type FdSet struct {
//     Bits [32]int32 // Max FD = 1024 = 32x32
// }

// FDZero set to zero the fdSet
func FDZero(p *syscall.FdSet) {
	p.Bits = [32]int32{}
}

// FDSet actives a given bit of fdSet
func FDSet(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] |= (1 << (uint(fd) % 32))
}

// FDClr actives a given bit of fdSet
func FDClr(fd int, p *syscall.FdSet) {
	p.Bits[fd/32] &^= (1 << (uint(fd) % 32))
}

// FDIsSet return true if the given fd is set
func FDIsSet(fd int, p *syscall.FdSet) bool {
	return p.Bits[fd/32]&(1<<(uint(fd)%32)) != 0
}
