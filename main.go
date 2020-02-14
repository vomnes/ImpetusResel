package main

import (
	"fmt"
	"log"
	"strconv"
	"syscall"
)

// socket, accept, listen, send, recv, bind, connect, inet_addr,
// setsockopt, getsockname

const (
	listenBacklog = 100
)

func main() {
	fmt.Println("Welcome in ImpetusResel")
	// AF_INET  0x2 -> The Internet Protocol version 4 (IPv4) address family
	// AF_INET6 0x1E -> The Internet Protocol version 6 (IPv6) address family
	// Socket types
	// SOCK_STREAM	1		     Stream (connection) socket for reliable, sequenced, connection oriented messages (think TCP)
	// SOCK_DGRAM	  2		     Datagram (conn.less) socket for connection-less, unreliable messages (think UDP or UNIX connections)
	// SOCK_RAW	    3		     Raw socket
	socketFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_IP)
	if err != nil {
		log.Fatalln("Socket -", err)
	}
	server := &syscall.SockaddrInet4{
		Port: 8084,
		Addr: [4]byte{127, 0, 0, 1},
	}
	err = syscall.Bind(socketFD, server)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to bind to Addr: %v, Port: %d\nReason: %s", server.Addr, server.Port, err))
	}
	fmt.Printf("Server: Bound to addr: %v, port: %d\n", server.Addr, server.Port)
	// syscall.Listen(sockfd, backlog int) error
	// sockfd, a valid socket descriptor
	// backlog, an integer representing the number of pending connections that can be queued up at any one time.
	err = syscall.Listen(socketFD, listenBacklog)
	if err != nil {
		log.Fatalln("Listen -", err)
	}
	for {
		nfd, sa, err := syscall.Accept(socketFD)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Connection to", nfd, sa)
			content := "Hello World"
			contentLen := len(content)
			// func Sendmsg(destFD int, p, oob []byte, to Sockaddr, flags int) error
			// destFD is the destinataire file descriptor
			// p is the content of the message
			// oob is the Out Of Band data
			// to is the receiver socket address
			// flags is the bitwise OR of zero or more of the following flags :
			// MSG_CONFIRM, MSG_DONTROUTE, MSG_DONTWAIT, MSG_EOR, MSG_MORE, MSG_NOSIGNAL, MSG_OOB
			err = syscall.Sendmsg(
				nfd,
				[]byte("HTTP/1.1 200 OK\r\nStatus: 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: "+strconv.Itoa(contentLen)+"\r\n\r\n"+content),
				nil, sa, syscall.MSG_DONTWAIT)
			if err != nil {
				fmt.Println("Sendmsg", err)
			}
		}
	}
}

// // func SetsockoptInet4Addr(fd, level, opt int, value [4]byte) error
// // level argument specifies the protocol level at which the option resides
// // option_name argument specifies a single option to set. The option_name argument and any specified options are passed uninterpreted to the appropriate protocol module for interpretations
// err = syscall.SetsockoptInet4Addr(socketFD, syscall.IPPROTO_IP, syscall.SO_REUSEADDR, [4]byte{78, 238, 249, 32})
// if err != nil {
// 	log.Fatalln("SetsockoptInet4Addr", err)
// }
