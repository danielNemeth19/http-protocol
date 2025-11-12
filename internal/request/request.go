package request

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/danielNemeth19/http-protocol/internal/headers"
)

var endLine = []byte("\r\n")

const bufferSize = 8

var methods = []string{
	"GET",
	"POST",
	"DELETE",
	"PATCH",
	"PUT",
	"OPTIONS",
	"HEAD",
	"TRACE",
	"CONNECT",
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState int

const (
	initialized parseState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		line, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if line != nil {
			r.RequestLine = *line
			r.state = requestStateParsingHeaders
			return n, err
		}
		return 0, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if !done {
			return n, err
		}
		r.state = requestStateParsingBody
		return n, err
	case requestStateParsingBody:
		contentLength := r.Headers.Get("content-length")
		if contentLength == "" {
			r.state = requestStateDone
			return 0, nil
		} else {
			length, err := strconv.Atoi(contentLength)
			if err != nil {
				return 0, err
			}
			r.Body = append(r.Body, data...)
			if len(r.Body) == length {
				r.state = requestStateDone
			}
			if len(r.Body) > length {
				fmt.Println(len(r.Body), length)
				// r.state = requestStateDone
				return 0, fmt.Errorf("Length of body (%d) is greater then the Content-Length (%d)", len(r.Body), length)
			}
			return len(data), nil
		}
	}
	return 0, fmt.Errorf("Not sure what's going on")
}

func (r *Request) parse(data []byte) (int, error) {
	totalParsed := 0
	if r.state == requestStateDone {
		return 0, fmt.Errorf("Error: trying to read data in done state")
	}
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalParsed:])
		if err != nil {
			return 0, fmt.Errorf("Error during parsing: %s", err)
		}
		if n == 0 {
			break
		}
		totalParsed += n
	}
	return totalParsed, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	var reqLine RequestLine

	endLineSep := bytes.Index(b, endLine)
	if endLineSep == -1 {
		return nil, 0, nil
	}

	line := b[:endLineSep]
	parts := strings.Split(string(line), " ")
	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("Request line supposed to have three parts, got: %d\n", len(parts))
	}
	method, target, protocolPart := parts[0], parts[1], parts[2]
	if !slices.Contains(methods, method) {
		return nil, 0, fmt.Errorf("%s is not a valid method\n", method)
	}
	parts = strings.Split(protocolPart, "/")
	if len(parts) != 2 || parts[1] != "1.1" {
		return nil, 0, fmt.Errorf("HTTP Version is unsupported: %s\n", parts[1])
	}
	version := parts[1]
	reqLine.Method = method
	reqLine.RequestTarget = target
	reqLine.HttpVersion = version
	return &reqLine, len(line) + len(endLine), nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0
	buf := make([]byte, bufferSize)
	req := Request{state: initialized, Headers: headers.NewHeaders()}
	for req.state != requestStateDone {
		if readToIndex >= cap(buf) {
			newBuf := make([]byte, 2*cap(buf))
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.state = requestStateDone
			break
		} else if err != nil {
			return nil, err
		}
		readToIndex += n
		parsedBytes, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if parsedBytes != 0 {
			remainderBytes := readToIndex - parsedBytes
			copy(buf, buf[parsedBytes:readToIndex])
			readToIndex = remainderBytes
		}
	}
	return &req, nil
}
