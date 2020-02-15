package http

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"../utils"
	"github.com/kr/pretty"
)

// Header allows to store headers
type Header map[string][]string

// Request is the structure where the request extracted data is stored
type Request struct {
	Method string
	URL    string

	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	Header Header

	Body string

	ContentLength int64

	Host string

	Form     url.Values
	PostForm url.Values
}

// NewRequest init a new request structure
func NewRequest() *Request {
	return &Request{
		Header: Header{},
	}
}

// Print print the request structure
func (r *Request) Print() {
	pretty.Print(r)
}

const (
	segmentRequest = iota
	segmentHeaders
	segmentBody
)

func validMethod(method string) bool {
	validMethods := []string{
		"OPTIONS", // Section 9.2
		"GET",     // Section ^ + 0.1
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"TRACE",
		"CONNECT",
	}
	return utils.StringInArray(method, validMethods)
}

// ParseHTTPVersion parses a HTTP version string. (src: net/http/request.go)
// "HTTP/1.0" returns (1, 0, true).
func parseHTTPVersion(value string) (major, minor int, ok bool) {
	const Big = 1000000
	switch value {
	case "HTTP/1.1":
		return 1, 1, true
	case "HTTP/1.0":
		return 1, 0, true
	}
	if !strings.HasPrefix(value, "HTTP/") {
		return 0, 0, false
	}
	dot := strings.Index(value, ".")
	if dot < 0 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(value[5:dot])
	if err != nil || major < 0 || major > Big {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(value[dot+1:])
	if err != nil || minor < 0 || minor > Big {
		return 0, 0, false
	}
	return major, minor, true
}

// RequestParse extract and store the data from the request
// in the request structure
func (r *Request) RequestParse(headers string) error {
	listErrors := []string{}
	array := strings.Split(headers, "\r\n")
	segment := 0
	var ok bool

	for _, header := range array {
		fmt.Println(segment, header)
		switch segment {
		case segmentRequest:
			request := strings.Split(header, " ")
			if len(request) != 3 {
				listErrors = append(listErrors, "Invalid resquest number elements")
			}
			if !validMethod(request[0]) {
				listErrors = append(listErrors, request[0]+" is not a valid method")
			} else {
				r.Method = request[0]
			}
			r.URL = request[1]
			r.ProtoMajor, r.ProtoMinor, ok = parseHTTPVersion(request[2])
			if !ok {
				listErrors = append(listErrors, request[2]+" is not a valid HTTP version")
			} else {
				r.Proto = request[2]
			}
			segment = segmentHeaders
		case segmentHeaders:
			if header == "" {
				segment = segmentBody
				continue
			}
		case segmentBody:

		default:
			break
		}
	}
	return nil
}
