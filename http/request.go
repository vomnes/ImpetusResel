package http

import (
	"fmt"
	"strconv"
	"strings"

	"../utils"
	"github.com/kr/pretty"
)

// Header allows to store headers
type Header map[string][]string

// Value store the URL values from thes forms
type Values map[string][]string

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

	Form        Values
	HasForm     bool
	PostForm    Values
	HasPostForm bool
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

func (r *Request) addHeader(key string, values []string) {
	for _, value := range values {
		r.Header[key] = append(r.Header[key], value)
	}
}

// RequestParse extract and store the data from the request
// in the request structure
func (r *Request) RequestParse(headers string) error {
	listErrors := []string{}
	array := strings.Split(headers, "\r\n")

	parts := strings.Split(headers, "\r\n\r\n")
	fmt.Println(parts[0])
	fmt.Println("======")
	fmt.Println(parts[1])
	segment := 0
	var ok bool

	for _, header := range array {
		// fmt.Println(segment, header)
		switch segment {
		case segmentRequest:
			request := strings.Split(header, " ")
			if len(request) != 3 {
				listErrors = append(listErrors, "Invalid Protocol/URL/Version format")
				continue
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
			h := strings.Split(header, ": ")
			if len(h) != 2 {
				listErrors = append(listErrors, "Invalid header format")
				continue
			}
			if h[0] == string(ContentLength) {
				value, err := strconv.Atoi(h[1])
				if err == nil {
					r.ContentLength = int64(value)
				}
				continue
			}
			if h[0] == string(Host) {
				r.Host = h[1]
				continue
			}
			if headerName(h[0]) == ContentType {
				if strings.Contains(h[1], "application/x-www-form-urlencoded") {
					r.HasForm = true
				}
				if strings.Contains(h[1], "multipart/form-data") {
					r.HasPostForm = true
				}
			}
			headerValues := strings.Split(h[1], ",")
			r.addHeader(h[0], headerValues)
		case segmentBody:
			r.Body += header + "\n"
		default:
			break
		}
	}
	return nil
}
