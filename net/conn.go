package net

import "golang.org/x/sys/unix"

// Conn store a socket connection
type Conn struct {
	Fd   int
	Addr unix.Sockaddr
}

// Read store in buf the data received from a socket connection
func (c *Conn) Read(buf *[]byte) (int, error) {
	// * Recvfrom will read the client fd and store the data in msg
	// Do not forger to close the fd after
	sizeMsg, _, err := unix.Recvfrom(c.Fd, *buf, 0)
	if err != nil {
		return 0, err
	}
	return sizeMsg, nil
}

// Read send the buf data to a socket connection
func (c *Conn) Write(buf []byte) error {
	return unix.Sendmsg(
		c.Fd,
		buf,
		nil,
		c.Addr,
		unix.MSG_DONTWAIT)
}

// Close closes the fd of a socket connection
func (c *Conn) Close() error {
	return unix.Close(c.Fd)
}
