package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
	// "strings"
)

var methods = []string{
	"GET", "POST",
	"DELETE", "PATCH",
	"PUT", "OPTIONS",
	"HEAD", "TRACE", "CONNECT",
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

const endLine = "\r\n"

type parseState int

const (
	initalized parseState = iota
	done
)

func (r *Request) parse(data []byte) (int, error) {
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

	curr_parts := strings.Split(string(b), endLine)
	if len(curr_parts) != 2 {
		return nil, 0, nil
	}

	line := curr_parts[0]
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
	return &reqLine, 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0
	buf := make([]byte, 8)
	req := Request{state: initalized}
	for req.state != done {
		n, err := reader.Read(buf[readToIndex:])
		readToIndex += n
		if err == io.EOF {
			fmt.Printf("EOF HERE\n%v\n", req.state)
			return nil, err
		}
		if readToIndex >= len(buf) {
			newBuf := make([]byte, 2*cap(buf))
			copy(newBuf, buf)
			buf = newBuf
		}
		n, err = req.parse(buf)
		if err != nil {
			return nil, err
		}
	}
	return &req, nil
}
