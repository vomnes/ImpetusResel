package http

import (
	"strconv"
)

// Hypertext Transfer Protocol (HTTP/1.1): Message Syntax and Routing

type headerName string

const (
	AcceptCharset    headerName = "Accept-Charset"
	AcceptEncoding   headerName = "Accept-Encoding"
	AcceptLanguage   headerName = "Accept-Language"
	Allow            headerName = "Allow"
	Authorization    headerName = "Authorization"
	ContentEncoding  headerName = "Content-Encoding"
	ContentLanguage  headerName = "Content-Language"
	ContentLength    headerName = "Content-Length"
	ContentLocation  headerName = "Content-Location"
	ContentType      headerName = "Content-Type"
	Date             headerName = "Date"
	Host             headerName = "Host"
	LastModified     headerName = "Last-Modified"
	Location         headerName = "Location"
	Referer          headerName = "Referer"
	RetryAfter       headerName = "Retry-After"
	Server           headerName = "Server"
	TransferEncoding headerName = "Transfer-Encoding"
	UserAgent        headerName = "User-Agent"
	WWWAuthenticate  headerName = "WWW-Authenticate"
)

// Headers ...
type Headers struct {
	version    string
	statusCode int
	entities   map[string]string
	body       string
}

// NewHeader init the headers structure
func NewHeader() *Headers {
	return &Headers{
		entities: map[string]string{},
	}
}

// GetStatus return the status code and status text under a string format
func (h Headers) GetStatus() string {
	return strconv.Itoa(h.statusCode) + " " + StatusText(h.statusCode)
}

// AddEntity add a new header entity
func (h *Headers) AddEntity(key headerName, value string) {
	h.entities[string(key)] = value
}

// SetVersion set the version value in headers
func (h *Headers) SetVersion(version string) { h.version = version }

// SetStatusCode set the status code value in headers
func (h *Headers) SetStatusCode(code int) { h.statusCode = code }

// SetBody set the body content in headers
func (h *Headers) SetBody(content string) { h.body = content }

// ToByte return the headers data under []byte format
func (h Headers) ToByte() []byte {
	var header string
	// HTTP/1.1 200 OK\r\nStatus: 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\nContent-Length: "+strconv.Itoa(contentLen)+"\r\n\r\n"+content
	header += "HTTP/" + h.version + " " + h.GetStatus() + "\r\n"
	header += "Status: " + h.GetStatus() + "\r\n"
	for key, value := range h.entities {
		header += key + ": " + value + "\r\n"
	}
	header += string(ContentLength) + ": " + strconv.Itoa(len(h.body)) + "\r\n\r\n"
	header += string(h.body)
	return []byte(header)
}
