package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

func isValidMethod(method string) bool {
	var methods = []string{"GET", "POST", "DELETE", "PATCH", "PUT", "OPTIONS", "HEAD", "TRACE", "CONNECT"}
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
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
		return nil, fmt.Errorf("Request line supposed to have three parts, got: %s\n", parts)
	}
	method, target, versionPart := parts[0], parts[1], parts[2]
	if !isValidMethod(method) {
		return nil, fmt.Errorf("%s is not a valid method\n", method)
	}
	vp := strings.Split(versionPart, "/")
	if len(vp) != 2 || vp[1] != "1.1" {
		return nil, fmt.Errorf("HTTP Version is unsupported: %s\n", versionPart)
	}
	reqLine.Method = method
	reqLine.RequestTarget = target
	reqLine.HttpVersion = vp[1]
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
