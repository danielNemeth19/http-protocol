package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

const endLine = "\r\n"

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
	done
)

func (r *Request) parse(data []byte) (int, error) {
	if r.state == done {
		return 0, fmt.Errorf("Error: trying to read data in done state")
	}
	line, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if line != nil {
		r.RequestLine = *line
		r.state = done
		return n, err
	}
	return 0, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	var reqLine RequestLine

	data := strings.Split(string(b), endLine)
	if len(data) != 2 {
		return nil, 0, nil
	}

	line := data[0]
	parts := strings.Split(line, " ")
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
	req := Request{state: initialized}
	for req.state != done {
		if readToIndex >= cap(buf) {
			newBuf := make([]byte, 2*cap(buf))
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}
		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			req.state = done
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
