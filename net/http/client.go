package http

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"

	"golang.org/x/sys/unix"
)

func NewRequest(method, url string, body []byte) (Request, error) {
	if method == "" {
		method = "GET"
	}
	if !validMethod(method) {
		return Request{}, errors.New("Invalid method")
	}
	host := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`).Find([]byte(url))
	request := Request{
		Method: method,
		Host:   string(host),
		URL:    url[len(host):],

		Proto: "HTTP/1.1",

		Header: Header{},

		ContentLength: int64(len(body)),
		Body:          []byte(body),
	}
	return request, nil
}

func Do(req *Request) (Request, error) {
	serverFD, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_IP)
	if err != nil {
		log.Fatal("Socket: ", err)
	}
	// Resolve IP Address DNS - Original
	hostIP, err := net.LookupIP(req.Host)
	if err != nil {
		fmt.Println(err)
		return Request{}, err
	}
	var addr net.IP
	for _, ip := range hostIP {
		if addr = ip.To4(); addr != nil {
			break
		}
	}
	serverAddr := &unix.SockaddrInet4{
		Port: 8085,
		Addr: [4]byte{addr[0], addr[1], addr[2], addr[3]},
	}

	fmt.Println("a")
	err = unix.Connect(serverFD, serverAddr)
	if err != nil {
		if err == unix.ECONNREFUSED {
			fmt.Println("* Connection failed")
		} else {
			fmt.Println("Connect:", err)
		}
		unix.Close(serverFD)
		return Request{}, err
	}
	fmt.Println("r")
	err = unix.Sendmsg(
		serverFD,
		req.Bytes(),
		nil, serverAddr, unix.MSG_DONTWAIT)
	if err != nil {
		fmt.Println("Sendmsg: ", err)
		return Request{}, err
	}

	buf := make([]byte, 8000)
	fmt.Println("y")
	_, _, err = unix.Recvfrom(serverFD, buf, 0)
	if err != nil {
		fmt.Println("Recvfrom: ", err)
		unix.Close(serverFD)
		return Request{}, err
	}
	fmt.Println("é", string(buf), "|")
	unix.Close(serverFD)
	resp := InitRequest()
	// err = resp.RequestParse(string(buf))
	// if err != nil {
	// 	fmt.Println("RequestParse: ", err)
	// 	return Request{}, err
	// }
	return *resp, nil
}
