package net

import (
	"syscall"
)

// type FdSet struct {
//     Bits [32]int32 // FD_SETSIZE = 1024 = 32x32
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

// FDAddr is the type storing the sockaddr of each fd
type FDAddr map[int]syscall.Sockaddr

// FDAddrInit init FDAddr with the size of FDSize
func FDAddrInit() *FDAddr {
	f := make(FDAddr, syscall.FD_SETSIZE)
	return &f
}

// Get return the Sockaddr value of a given fd key
func (f *FDAddr) Get(fd int) syscall.Sockaddr {
	return (*f)[fd]
}

// Set set the Sockaddr value of a given fd key
func (f *FDAddr) Set(fd int, value syscall.Sockaddr) {
	(*f)[fd] = value
}

// Clr remove a given fd key in FDAddr
func (f *FDAddr) Clr(fd int) {
	delete(*f, fd)
}
