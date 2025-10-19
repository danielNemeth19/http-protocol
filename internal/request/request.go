package request

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

var methods = []string{
	"GET", "POST",
	"DELETE", "PATCH",
	"PUT", "OPTIONS",
	"HEAD", "TRACE", "CONNECT",
}

type Request struct {
	RequestLine RequestLine
	initalized int
	done int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func parseRequestLine(line string) (*RequestLine, error) {
	var reqLine RequestLine
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Request line supposed to have three parts, got: %d\n", len(parts))
	}
	method, target, protocolPart := parts[0], parts[1], parts[2]
	if !slices.Contains(methods, method) {
		return nil, fmt.Errorf("%s is not a valid method\n", method)
	}
	parts = strings.Split(protocolPart, "/")
	if len(parts) != 2 || parts[1] != "1.1" {
		return nil, fmt.Errorf("HTTP Version is unsupported: %s\n", parts[1])
	}
	version := parts[1]
	reqLine.Method = method
	reqLine.RequestTarget = target
	reqLine.HttpVersion = version
	return &reqLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	var request Request
	i, err := io.ReadAll(reader)
	if err != nil {
		panic("Input cannot be read")
	}
	parts := strings.Split(string(i), "\r\n")
	requestLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}
	request.RequestLine = *requestLine
	return &request, nil
}
