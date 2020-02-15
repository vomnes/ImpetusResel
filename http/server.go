package http

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"../utils"
)

// socket, accept, listen, send, recv, bind, connect, inet_addr,
// setsockopt, getsockname

const (
	listenBacklog = 100
)

type server struct {
	fd int
}

type data struct {
	server server
	router *Router
}

func (s *data) Socket() error {
	var err error
	// AF_INET  0x2 -> The Internet Protocol version 4 (IPv4) address family
	// AF_INET6 0x1E -> The Internet Protocol version 6 (IPv6) address family
	// Socket types
	// SOCK_STREAM	1		     Stream (connection) socket for reliable, sequenced, connection oriented messages (think TCP)
	// SOCK_DGRAM	  2		     Datagram (conn.less) socket for connection-less, unreliable messages (think UDP or UNIX connections)
	// SOCK_RAW	    3		     Raw socket
	s.server.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)
	return err
}

func (s *data) Listen() error {
	// syscall.Listen(sockfd, backlog int) error
	// sockfd, a valid socket descriptor
	// backlog, an integer representing the number of pending connections that can be queued up at any one time.
	return syscall.Listen(s.server.fd, listenBacklog)
}

func (s *data) SetRouter(router *Router) {
	s.router = router
}

// Client ...
type Client struct {
	fd        int
	stockaddr syscall.Sockaddr
}

func (s *data) Send(h *Headers, dst Client) error {
	// func Sendmsg(destFD int, p, oob []byte, to Sockaddr, flags int) error
	// destFD is the destinataire file descriptor
	// p is the content of the message
	// oob is the Out Of Band data
	// to is the receiver socket address
	// flags is the bitwise OR of zero or more of the following flags :
	// MSG_CONFIRM, MSG_DONTROUTE, MSG_DONTWAIT, MSG_EOR, MSG_MORE, MSG_NOSIGNAL, MSG_OOB
	return syscall.Sendmsg(
		dst.fd,
		h.ToByte(),
		nil, dst.stockaddr, syscall.MSG_DONTWAIT)
}

func (s *data) handleRoute() {
	nfd, sa, err := syscall.Accept(s.server.fd)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connection to", nfd, sa)
		c := Client{
			fd:        nfd,
			stockaddr: sa,
		}
		file := os.NewFile(uintptr(c.fd), "")
		if file == nil {
			fmt.Println("NewFile", err)
			return
		}
		data := make([]byte, 8000)
		_, err := file.Read(data)
		if err != nil {
			log.Fatal(err)
		}

		h := NewHeader()
		h.SetVersion("1.1")
		r := NewRequest()
		r.RequestParse(string(data))
		route := s.router.routes[r.URL]
		if route.Handler != nil {
			route.Handler(h, r)
		} else {
			s.router.defaultHandler(h, r)
		}
		defer file.Close()
		err = s.Send(h, c)
		if err != nil {
			fmt.Println("Send", err)
		}
	}
}

// ListenAndServe will launch the server on a given port
func ListenAndServe(port int, router *Router) {
	n := data{}
	n.SetRouter(router)
	err := n.Socket()
	if err != nil {
		log.Fatalln("Socket -", err)
	}
	server := &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{127, 0, 0, 1},
	}
	err = syscall.Bind(n.server.fd, server)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to bind to Addr: %v, Port: %d\nReason: %s", utils.ByteArrayJoin(server.Addr[:], "."), server.Port, err))
	}
	fmt.Printf("Server: Bound to addr: %v, port: %d\n", utils.ByteArrayJoin(server.Addr[:], "."), server.Port)
	err = n.Listen()
	if err != nil {
		log.Fatalln("Listen -", err)
	}
	for {
		n.handleRoute()
	}
}

// // func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) error
// // level argument specifies the protocol level at which the option resides
// // option_name argument specifies a single option to set. The option_name argument and any specified options are passed uninterpreted to the appropriate protocol module for interpretations
// err = syscall.SetsockoptInet4Addr(socketFD, syscall.IPPROTO_IP, syscall.SO_REUSEADDR, [4]byte{78, 238, 249, 32})
// if err != nil {
// 	log.Fatalln("SetsockoptInet4Addr", err)
// }
