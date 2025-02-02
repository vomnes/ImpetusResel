package http

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"../../utils"
	"github.com/kylelemons/godebug/pretty"
)

var (
	METHODS = []string{
		"OPTIONS", // Section 9.2
		"GET",     // Section ^ + 0.1
		"HEAD",
		"POST",
		"PUT",
		"DELETE",
		"TRACE",
		"CONNECT",
	}
)

// URL Living Standard - https://url.spec.whatwg.org/#urlencoded-parsing

// Header allows to store headers
type Header map[string][]string

func (h *Header) AddHeader(key string, value string) {
	(*h)[key] = append((*h)[key], value)
}

func (h *Header) AddHeaders(key string, values []string) {
	for _, value := range values {
		h.AddHeader(key, strings.TrimSpace(value))
	}
}

// IsSet return true if the key has at least a value
func (h Header) IsSet(key string) bool {
	return len(h[key]) != 0
}

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

	Body []byte

	ContentLength int64

	Host string

	Form        Values
	HasForm     bool
	PostForm    Values
	HasPostForm bool

	ParsingError []string
}

// InitRequest init a new request structure
func InitRequest() *Request {
	return &Request{
		Header:   Header{},
		Form:     Values{},
		PostForm: Values{},
	}
}

// Bytes returns the request under HTTP format in []byte
func (r *Request) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteString(r.Method + " " + r.URL + " " + r.Proto + "\r\n")
	buf.WriteString("Host: " + r.Host + "\r\n")
	buf.WriteString("User-Agent: Go\r\n")
	buf.WriteString("Accept: */*\r\n")
	for key, values := range r.Header {
		buf.WriteString(key + ": ")
		for i, value := range values {
			buf.WriteString(value + ",")
			if i < len(values)-1 {
				buf.WriteString(value + ",")
			}
		}
		buf.WriteString("\r\n")
	}
	buf.WriteString(string(ContentLength) + ": " + strconv.Itoa(len(r.Body)) + "\r\n\r\n")
	buf.Write(r.Body)
	return buf.Bytes()
}

// > GET /books/v1/volumes?q=isbn:0747532699 HTTP/2
// > Host: www.googleapis.com
// > User-Agent: curl/7.54.0
// > Accept: */*

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
	return utils.StringInArray(method, METHODS)
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
			r.Header.AddHeaders(h[0], headerValues)
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

func (r *Request) getMultipartBoundaryDelimiter() string {
	contentType := r.Header[string(ContentType)]
	if contentType == nil {
		return ""
	}
	subpart := "boundary="
	multipart := strings.LastIndex(contentType[0], subpart)
	if multipart != -1 {
		return strings.ReplaceAll(contentType[0][multipart+len(subpart):], "\"", "")
	}
	return ""
}

func extractDataFormContent(content string) string {
	content = strings.ReplaceAll(content, "\r\n", " ")
	content = strings.ReplaceAll(content, ":", "")
	content = strings.ReplaceAll(content, ";", "")
	content = strings.ReplaceAll(content, "\"", "")
	arr := strings.Split(content, " ")

	var lenItem int
	for _, item := range arr {
		lenItem = len(item)
		// if strings.Contains(item, "filename=") {
		// if lenItem > len("filename=") {
		// 	name += item[len("filename="):] + ";"
		// }
		// } else
		if strings.Contains(item, "name=") {
			if lenItem > len("name=") {
				return item[len("name="):]
			}
		}
		// if strings.Contains(item, "/") {
		// 	name += item
		// }
	}
	return ""
}

func (r *Request) parseDataForm(body string) {
	var parts, items []string
	var partSize int
	// var nameItem string

	delimiter := "--" + r.getMultipartBoundaryDelimiter()
	parts = strings.Split(body, delimiter)
	for _, part := range parts {
		partSize = len(part)
		// Skip part where "Content-Disposition" is not include
		if partSize < len("Content-Disposition") {
			continue
		}
		part = strings.TrimLeft(part, "\r\n")
		// Handle part
		items = strings.Split(part, "\r\n\r\n")
		if len(items) == 2 {
			r.PostForm.AddValue(
				extractDataFormContent(items[0]),
				items[1],
			)
		}
	}
	return
}

func (r *Request) parseBody(body string) {
	r.Body = []byte(body)
	if r.HasForm {
		r.parseForm(body)
	}
	if r.HasPostForm {
		r.parseDataForm(body)
	}
}

// RequestParse extract and store the data from the request
// in the request structure
func (r *Request) RequestParse(headers string) error {
	const delimiter = "\r\n\r\n"
	bodyStart := strings.Index(headers, delimiter)
	if bodyStart == -1 {
		return errors.New("Not a valid reader format")
	}
	r.parseHeaders(headers[:bodyStart])
	r.parseBody(headers[bodyStart+len(delimiter):])
	return nil
}
