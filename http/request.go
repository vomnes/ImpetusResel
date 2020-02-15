package http

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"../utils"
	"github.com/kr/pretty"
)

// URL Living Standard - https://url.spec.whatwg.org/#urlencoded-parsing

// Header allows to store headers
type Header map[string][]string

// Values store the URL values from thes forms
type Values map[string][]string

// AddValue append a new item to a given key in Values
func (v *Values) AddValue(key, value string) {
	(*v)[key] = append((*v)[key], value)
}

// Request is the structure where the request extracted data is stored
type Request struct {
	Method string
	URL    string

	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	Header Header

	Body io.Reader

	ContentLength int64

	Host string

	Form        Values
	HasForm     bool
	PostForm    Values
	HasPostForm bool

	ParsingError []string
}

// NewRequest init a new request structure
func NewRequest() *Request {
	return &Request{
		Header:   Header{},
		Form:     Values{},
		PostForm: Values{},
	}
}

// Print print the request structure
func (r *Request) Print() {
	pretty.Print(r)
}

func (r *Request) pushError(err string) {
	r.ParsingError = append(r.ParsingError, err)
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

func (r *Request) parseHeaders(headers string) {
	array := strings.Split(headers, "\r\n")
	var ok, requestSpecCollected bool

	for _, header := range array {
		if !requestSpecCollected {
			request := strings.Split(header, " ")
			if len(request) != 3 {
				r.pushError("Invalid Protocol/URL/Version format")
				continue
			}
			if !validMethod(request[0]) {
				r.pushError(request[0] + " is not a valid method")
			} else {
				r.Method = request[0]
			}
			r.URL = request[1]
			r.ProtoMajor, r.ProtoMinor, ok = parseHTTPVersion(request[2])
			if !ok {
				r.pushError(request[2] + " is not a valid HTTP version")
			} else {
				r.Proto = request[2]
			}
			requestSpecCollected = true
		} else {
			h := strings.Split(header, ": ")
			if len(h) != 2 {
				r.pushError("Invalid header format")
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
		}
	}
}

// https://url.spec.whatwg.org/#urlencoded-parsing
// application/x-www-form-urlencoded parsing
func (r *Request) parseForm(form string) {
	var nameString, valueString string
	urlValues := strings.Split(form, "&")
	for _, value := range urlValues {
		keyValue := strings.Split(value, "=")
		if len(keyValue) == 2 {
			nameString = keyValue[0]
			valueString = keyValue[1]
			// Replace any 0x2B (+) in name and value with 0x20 (SP).
			valueString = strings.ReplaceAll(valueString, "+", " ")
			r.Form.AddValue(nameString, valueString)
		}
	}
}

func (r *Request) parseBody(body string) {
	r.Body = strings.NewReader(body)
	if r.HasForm {
		r.parseForm(body)
	}
}

// RequestParse extract and store the data from the request
// in the request structure
func (r *Request) RequestParse(headers string) error {
	fmt.Println(headers)
	fmt.Println("=============================================================")

	segements := strings.Split(headers, "\r\n\r\n")
	r.parseHeaders(segements[0])
	if r.HasPostForm {
		r.parseBody(headers)
	}

	fmt.Println(segements)
	return nil
}
