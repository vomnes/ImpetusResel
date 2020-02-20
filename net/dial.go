package net

import (
	"fmt"

	"golang.org/x/sys/unix"

	"../utils"
)

const (
	listenBacklog = 100
)

type TCPServer struct {
	Fd       int
	AddrIPv4 *unix.SockaddrInet4
}

func initSockAddr(addr string, port int) *unix.SockaddrInet4 {
	tmpAddr := ParseIP(addr)
	if tmpAddr == nil {
		tmpAddr = IP{127, 0, 0, 1}
	}
	return &unix.SockaddrInet4{
		Port: port,
		Addr: [4]byte{tmpAddr[0], tmpAddr[1], tmpAddr[2], tmpAddr[3]},
	}
}

func initTCPServer(addr string, port int) TCPServer {
	return TCPServer{
		Fd:       0,
		AddrIPv4: initSockAddr(addr, port),
	}
}

func (s *TCPServer) socket() error {
	var err error
	// * Socket will return the server socket file descriptor
	s.Fd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)
	return err
}

func (s *TCPServer) Listen() error {
	// * Listen will set sockfd as a passive socket ready to accept
	// incoming connection request
	return unix.Listen(s.Fd, listenBacklog)
}

// Dial creates the TCP connection, link the given address and port
// and start to listen
func Dial(port int) (TCPServer, error) {
	s := initTCPServer("127.0.0.1", port)
	err := s.socket()
	if err != nil {
		return TCPServer{}, fmt.Errorf("socket: %s", err.Error())
	}
	// * Bind will link a socket file descriptor to a socket address
	err = unix.Bind(s.Fd, s.AddrIPv4)
	if err != nil {
		return TCPServer{}, fmt.Errorf("Failed to bind to Addr: %v, Port: %d\nReason: %s", utils.ByteArrayJoin(s.AddrIPv4.Addr[:], "."), s.AddrIPv4.Port, err.Error())
	}
	return s, nil
}

// Accept accepts a connection on the TCPServer and return this connection
func (s *TCPServer) Accept() (Conn, error) {
	// * Accept extracts the first connection request on the queue of
	// pending connections for the listening socket, sockfd, creates a new
	// connected socket, and returns a new file descriptor referring
	// to that socket and the address of this socket.
	connFd, connAddr, err := unix.Accept(s.Fd)
	if err != nil {
		return Conn{}, err
	}
	return Conn{
		Fd:   connFd,
		Addr: connAddr,
	}, nil
}

// GetAddr returns a string formated containing the address and port
func (s *TCPServer) GetAddr() string {
	return fmt.Sprintf("%d:%d", s.AddrIPv4.Addr, s.AddrIPv4.Port)
}
