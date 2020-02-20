package http

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"

	"../../net"
	"../../utils"
)

// socket, accept, listen, send, recv, bind, connect, inet_addr,
// setsockopt, getsockname

// https://www.gnu.org/software/libc/manual/html_node/Sockets.html#Sockets
// https://www.gnu.org/software/libc/manual/html_node/Connections.html
// https://www.gnu.org/software/libc/manual/html_node/Server-Example.html
// https://www.tenouk.com/Module41.html

const (
	listenBacklog = 100
)

type server struct {
	fd       int
	addrIPv4 *unix.SockaddrInet4
}

type data struct {
	server server
	router *Router
}

func (s *data) Socket() error {
	var err error
	// * Socket will return the server socket file descriptor
	s.server.fd, err = unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)
	return err
}

func (s *data) Listen() error {
	// * Listen will set sockfd as a passive socket ready to accept
	// incoming connection request
	return unix.Listen(s.server.fd, listenBacklog)
}

func (s *data) SetRouter(router *Router) {
	s.router = router
}

func initSockAddr(addr string, port int) *unix.SockaddrInet4 {
	tmpAddr := net.ParseIP(addr)
	if tmpAddr == nil {
		tmpAddr = net.IP{127, 0, 0, 1}
	}
	return &unix.SockaddrInet4{
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
	stockaddr unix.Sockaddr
}

func (s *data) Send(h *Headers, fdClient int, addrClient unix.Sockaddr) error {
	// * Sendmsg will send a message on the socket connection
	return unix.Sendmsg(
		fdClient,
		h.ToByte(),
		nil, addrClient, unix.MSG_DONTWAIT)
}

var activeFdSet unix.FdSet

func (s *data) run() {
	var tmpFdSet unix.FdSet

	net.FDZero(&activeFdSet)
	net.FDSet(s.server.fd, &activeFdSet)
	fdAddr := net.FDAddrInit()

	for {
		tmpFdSet = activeFdSet

		// * Select will disable in the FdSet copy the FDs that
		// are not yet ready to be read
		err := unix.Select(unix.FD_SETSIZE, &tmpFdSet, nil, nil, nil)
		if err != nil {
			log.Fatal("Select ", err)
		}
		for fd := 0; fd < unix.FD_SETSIZE; fd++ {
			if net.FDIsSet(fd, &tmpFdSet) {
				if fd == s.server.fd {
					// * Accept extracts the first connection request on the queue of
					// pending connections for the listening socket, sockfd, creates a new
					// connected socket, and returns a new file descriptor referring
					// to that socket and the address of this socket.
					newFD, sa, err := unix.Accept(s.server.fd)
					if err != nil {
						fmt.Println("Accept", err)
						return
					}
					// == Add new connection fd == //
					net.FDSet(newFD, &activeFdSet)
					fdAddr.Set(newFD, sa)
				} else {
					msg := make([]byte, 1024)
					// * Recvfrom will read the client fd and store the data in msg
					// Do not forger to close the fd after
					sizeMsg, _, err := unix.Recvfrom(fd, msg, 0)
					if err != nil {
						net.FDClr(fd, &activeFdSet)
						fdAddr.Clr(fd)
						unix.Close(fd)
						continue
					}
					saFrom := fdAddr.Get(fd).(*unix.SockaddrInet4)
					fmt.Printf("%d byte read from %d:%d on socket %d\n",
						sizeMsg, saFrom.Addr, saFrom.Port, fd)
					// == Parse recv message - HTTP Type == //
					h := NewHeader()
					h.SetVersion("1.1")
					r := NewRequest()
					r.RequestParse(string(msg))
					fmt.Println("Message:", r.Method, r.URL)
					route := s.router.routes[r.URL]
					if route.Handler != nil {
						route.Handler(h, r)
					} else {
						s.router.defaultHandler(h, r)
					}
					// * Sendmsg will send a message on the socket connection
					err = s.Send(h, fd, fdAddr.Get(fd))
					if err != nil {
						fmt.Println("Send", err)
					}
					net.FDClr(fd, &activeFdSet)
					fdAddr.Clr(fd)
					unix.Close(fd)
				}
			}
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
	n.SetSocketAddr("127.0.0.1", port)
	// * Bind will link a socket file descriptor to a socket address
	err = unix.Bind(n.server.fd, n.server.addrIPv4)
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
// err = unix.SetsockoptInet4Addr(socketFD, unix.IPPROTO_IP, unix.SO_REUSEADDR, [4]byte{78, 238, 249, 32})
// if err != nil {
// 	log.Fatalln("SetsockoptInet4Addr", err)
// }
