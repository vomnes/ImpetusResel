package http

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestRequestParseHeaderForm(t *testing.T) {
	r := InitRequest()
	r.RequestParse("GET /form HTTP/1.1\r\nHost: localhost:8084\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\r\nContent-Length: 42\r\nContent-Type: application/x-www-form-urlencoded\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\nfield1=value1&field2=value2&field2=Hello+world !")
	diff := pretty.Compare(&r, Request{
		Method: "GET",
		Header: map[string][]string{
			"Accept": []string{
				"text/html",
				"application/xhtml+xml",
				"application/xml;q=0.9",
				"image/webp",
				"*/*;q=0.8",
			},
			"Accept-Encoding": []string{
				"gzip",
				"deflate",
			},
			"Accept-Language": []string{
				"en-US",
				"en;q=0.5",
			},
			"Connection":                []string{"keep-alive"},
			"Content-Type":              []string{"application/x-www-form-urlencoded"},
			"Upgrade-Insecure-Requests": []string{"1"},
			"User-Agent":                []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0"},
		},
		Body:          []byte("field1=value1&field2=value2&field2=Hello+world !"),
		ContentLength: 42,
		Host:          "localhost:8084",
		Form: Values{
			"field1": []string{
				"value1",
			},
			"field2": []string{
				"value2",
				"Hello world !",
			},
		},
		HasForm:      true,
		PostForm:     Values{},
		HasPostForm:  false,
		ParsingError: []string{},
		Proto:        "HTTP/1.1",
		ProtoMajor:   1,
		ProtoMinor:   1,
		URL:          "/form",
	})
	if diff != "" {
		t.Error(diff)
	}
}

func TestRequestParseHeaderMultipart(t *testing.T) {
	r := InitRequest()
	r.RequestParse("POST / HTTP/1.1\r\nHost: localhost:8084\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-US,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nContent-Type: multipart/form-data; boundary=\"---------------------------20762440193078419611623191500\"\r\nContent-Length: 772\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nvalentin omnes\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file1\"; filename=\"\"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file2\"; filename=\"http.html\"\r\nContent-Type: text/html\r\n\r\n<form action=\"http://localhost:8084\" method=\"post\" enctype=\"multipart/form-data\">\n  <p><input type=\"text\" name=\"text\" value=\"text default\">\n  <p><input type=\"file\" name=\"file1\">\n  <p><input type=\"file\" name=\"file2\">\n  <p><button type=\"submit\">Submit</button>\n</form>\n\r\n-----------------------------20762440193078419611623191500--\r\n")
	diff := pretty.Compare(&r, Request{
		Method: "POST",
		Header: map[string][]string{
			"Accept": []string{
				"text/html",
				"application/xhtml+xml",
				"application/xml;q=0.9",
				"image/webp",
				"*/*;q=0.8",
			},
			"Accept-Encoding": []string{
				"gzip",
				"deflate",
			},
			"Accept-Language": []string{
				"en-US",
				"en;q=0.5",
			},
			"Connection":                []string{"keep-alive"},
			"Content-Type":              []string{"multipart/form-data; boundary=\"---------------------------20762440193078419611623191500\""},
			"Upgrade-Insecure-Requests": []string{"1"},
			"User-Agent":                []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:74.0) Gecko/20100101 Firefox/74.0"},
		},
		Body:          []byte("-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nvalentin omnes\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file1\"; filename=\"\"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n-----------------------------20762440193078419611623191500\r\nContent-Disposition: form-data; name=\"file2\"; filename=\"http.html\"\r\nContent-Type: text/html\r\n\r\n<form action=\"http://localhost:8084\" method=\"post\" enctype=\"multipart/form-data\">\n  <p><input type=\"text\" name=\"text\" value=\"text default\">\n  <p><input type=\"file\" name=\"file1\">\n  <p><input type=\"file\" name=\"file2\">\n  <p><button type=\"submit\">Submit</button>\n</form>\n\r\n-----------------------------20762440193078419611623191500--\r\n"),
		ContentLength: 772,
		Host:          "localhost:8084",
		Form:          Values{},
		HasForm:       false,
		PostForm: map[string][]string{"file1": []string{"\r\n"},
			"file2": []string{"<form action=\"http://localhost:8084\" method=\"post\" enctype=\"multipart/form-data\">\n  <p><input type=\"text\" name=\"text\" value=\"text default\">\n  <p><input type=\"file\" name=\"file1\">\n  <p><input type=\"file\" name=\"file2\">\n  <p><button type=\"submit\">Submit</button>\n</form>\n\r\n"},
			"text":  []string{"valentin omnes\r\n"},
		},
		HasPostForm:  true,
		ParsingError: []string{},
		Proto:        "HTTP/1.1",
		ProtoMajor:   1,
		ProtoMinor:   1,
		URL:          "/",
	})
	if diff != "" {
		t.Error(diff)
	}
}
