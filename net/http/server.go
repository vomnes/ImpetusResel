package http

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"../../net"
	"../../utils"
)

// socket, accept, listen, send, recv, bind, connect, inet_addr,
// setsockopt, getsockname

const (
	listenBacklog = 100
)

type server struct {
	fd       int
	addrIPv4 *syscall.SockaddrInet4
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

func initSockAddr(addr string, port int) *syscall.SockaddrInet4 {
	tmpAddr := net.ParseIP(addr)
	if tmpAddr == nil {
		tmpAddr = net.IP{127, 0, 0, 1}
	}
	return &syscall.SockaddrInet4{
		Port: port,
		Addr: [4]byte{tmpAddr[0], tmpAddr[1], tmpAddr[2], tmpAddr[3]},
	}
}

func (s *data) SetSocketAddr(addr string, port int) {
	s.server.addrIPv4 = initSockAddr(addr, port)
}

// Client ...
type Client struct {
	fd        int
	stockaddr syscall.Sockaddr
}

func (s *data) Send(h *Headers, fdClient int, addrClient syscall.Sockaddr) error {
	// func Sendmsg(destFD int, p, oob []byte, to Sockaddr, flags int) error
	// destFD is the destinataire file descriptor
	// p is the content of the message
	// oob is the Out Of Band data
	// to is the receiver socket address
	// flags is the bitwise OR of zero or more of the following flags :
	// MSG_CONFIRM, MSG_DONTROUTE, MSG_DONTWAIT, MSG_EOR, MSG_MORE, MSG_NOSIGNAL, MSG_OOB
	return syscall.Sendmsg(
		fdClient,
		h.ToByte(),
		nil, addrClient, syscall.MSG_DONTWAIT)
}

func readFromClient(fd int) (msg []byte, file *os.File, err error) {
	fmt.Println("§-§")
	file = os.NewFile(uintptr(fd), "")
	if file == nil {
		err = errors.New("Not a valid file descriptor")
		return
	}
	fmt.Println("*-*")
	msg = make([]byte, 1024)
	n, err := file.Read(msg)
	fmt.Println(n, "EOF")
	fmt.Println("$-$")
	return
}

var activeFdSet syscall.FdSet

func (s *data) run() {
	var readFdSet syscall.FdSet

	net.FDZero(&activeFdSet)

	fmt.Println(s.server.fd, "run()")

	fdAddr := net.FDAddrInit()

	for {
		readFdSet = activeFdSet
		net.FDSet(s.server.fd, &activeFdSet)

		fmt.Println("readFdSet", readFdSet)
		fmt.Println("activeFdSet", activeFdSet)

		// func Select(int nfds, fd_set *FdSet, fd_set *FdSet, fd_set *FdSet, timeval *Timeval) error
		// -> ndfs : The select function checks only the first nfds file descriptors.
		// The usual thing is to pass FD_SETSIZE as the value of this argument.
		// -> fd_set : Data type represents file descriptor sets for the select function
		// -> timeval : The timeout specifies the maximum time to wait. If you pass
		// a null pointer for this argument, it means to block indefinitely until
		// one of the file descriptors is ready.
		// Specify zero as the time (a struct timeval containing all zeros)
		// if you want to find out which descriptors are ready without waiting if none are ready.
		// var timeval = syscall.Timeval{
		// 	Sec:  0,
		// 	Usec: 0,
		// }
		err := syscall.Select(syscall.FD_SETSIZE, &activeFdSet, nil, nil, nil)
		if err != nil {
			log.Fatal("Select ", err)
		}
		fmt.Println("* after:")
		fmt.Println("readFdSet", readFdSet)
		fmt.Println("activeFdSet", activeFdSet)
		for fd := 0; fd < syscall.FD_SETSIZE; fd++ {
			if net.FDIsSet(fd, &readFdSet) {
				if fd == s.server.fd {
					fmt.Println("a")
					fmt.Println("readFdSet", readFdSet)
					fmt.Println("activeFdSet", activeFdSet)
					newFD, sa, err := syscall.Accept(s.server.fd)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("-> readFdSet", readFdSet)
					fmt.Println("-> activeFdSet", activeFdSet)
					net.FDSet(newFD, &activeFdSet)
					fdAddr.Set(newFD, sa)
				} else {
					fmt.Println("$ Start", readFdSet)
					msg, _, _ := readFromClient(fd)
					if err != nil {
						fmt.Println("readFromClient", err)
						return
					}
					fmt.Println(string(msg))
					h := NewHeader()
					h.SetVersion("1.1")
					r := NewRequest()
					r.RequestParse(string(msg))
					route := s.router.routes[r.URL]
					if route.Handler != nil {
						route.Handler(h, r)
					} else {
						s.router.defaultHandler(h, r)
					}
					err = s.Send(h, fd, fdAddr.Get(fd))
					if err != nil {
						fmt.Println("Send", err)
					}
					// file.Close()
					net.FDClr(fd, &activeFdSet)
					fmt.Println("*> Done", readFdSet)
					fdAddr.Clr(fd)
					// return
				}
			}
		}
	}
}

// https://www.gnu.org/software/libc/manual/html_node/Sockets.html#Sockets
// https://www.gnu.org/software/libc/manual/html_node/Connections.html

// https://www.gnu.org/software/libc/manual/html_node/Server-Example.html

// ListenAndServe will launch the server on a given port
func ListenAndServe(port int, router *Router) {
	n := data{}
	n.SetRouter(router)
	err := n.Socket()
	if err != nil {
		log.Fatalln("Socket -", err)
	}
	n.SetSocketAddr("127.0.0.1", port)
	// Link socket to address IP
	err = syscall.Bind(n.server.fd, n.server.addrIPv4)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to bind to Addr: %v, Port: %d\nReason: %s", utils.ByteArrayJoin(n.server.addrIPv4.Addr[:], "."), n.server.addrIPv4.Port, err))
	}
	fmt.Printf("Server: Bound to addr: %v, port: %d\n", utils.ByteArrayJoin(n.server.addrIPv4.Addr[:], "."), n.server.addrIPv4.Port)
	err = n.Listen()
	if err != nil {
		log.Fatalln("Listen -", err)
	}
	n.run()
}

// // func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) error
// // level argument specifies the protocol level at which the option resides
// // option_name argument specifies a single option to set. The option_name argument and any specified options are passed uninterpreted to the appropriate protocol module for interpretations
// err = syscall.SetsockoptInet4Addr(socketFD, syscall.IPPROTO_IP, syscall.SO_REUSEADDR, [4]byte{78, 238, 249, 32})
// if err != nil {
// 	log.Fatalln("SetsockoptInet4Addr", err)
// }
