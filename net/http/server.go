package http

import (
	"fmt"

	"../../net"
)

// socket, accept, listen, send, recv, bind, connect, inet_addr,
// setsockopt, getsockname

// https://www.gnu.org/software/libc/manual/html_node/Sockets.html#Sockets
// https://www.gnu.org/software/libc/manual/html_node/Connections.html
// https://www.gnu.org/software/libc/manual/html_node/Server-Example.html
// https://www.tenouk.com/Module41.html

type server struct {
	socket net.TCPServer
	router *Router
}

func (s *server) SetRouter(router *Router) {
	s.router = router
}

func (s *server) run() {
	for {
		c, err := s.socket.Accept()
		if err != nil {
			fmt.Println("Accept:", err)
			continue
		}
		fmt.Println("Connection accepted on port:", c.Fd)
		go func(c net.Conn) {
			msg := make([]byte, 1024)
			_, err = c.Read(&msg)
			if err != nil {
				fmt.Println("Read:", err)
				c.Close()
			}

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
			err = c.Write(h.ToByte())
			if err != nil {
				fmt.Println("Write:", err)
				c.Close()
			}
			c.Close()
		}(c)
	}
}

// ListenAndServe will launch the server on a given port
func ListenAndServe(port int, router *Router) {
	s := server{}
	s.SetRouter(router)
	tcpSocket, err := net.Dial(port)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.socket = tcpSocket
	err = s.socket.Listen()
	if err != nil {
		fmt.Printf("Listen: %s", err.Error())
		return
	}
	fmt.Printf("Server is running on %s\n", s.socket.GetAddr())
	s.run()
}

// // func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) error
// // level argument specifies the protocol level at which the option resides
// // option_name argument specifies a single option to set. The option_name argument and any specified options are passed uninterpreted to the appropriate protocol module for interpretations
// err = unix.SetsockoptInet4Addr(socketFD, unix.IPPROTO_IP, unix.SO_REUSEADDR, [4]byte{78, 238, 249, 32})
// if err != nil {
// 	log.Fatalln("SetsockoptInet4Addr", err)
// }
