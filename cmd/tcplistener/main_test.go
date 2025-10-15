package main

import (
	"io"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

type HTTPMessage struct {
	StartLine  string
	FieldLines []string
	Body       string
}

func (h *HTTPMessage) expectedLines() []string {
	var m []string
	m = append(m, h.StartLine)
	for _, line := range h.FieldLines {
		m = append(m, line)
	}
	m = append(m, "\r\n")
	if h.Body != "" {
		m = append(m, h.Body)
	}
	return m
}

func (h *HTTPMessage) createHTTPMessage() string {
	endLine := "\r\n"
	var message strings.Builder
	message.WriteString(h.StartLine)
	message.WriteString(endLine)
	for _, headerRow := range h.FieldLines {
		message.WriteString(headerRow)
		message.WriteString(endLine)
	}
	message.WriteString(endLine)
	if h.Body != "" {
		message.WriteString(h.Body)
	}
	return message.String()
}

func (h *HTTPMessage) httpMessageAsReadCloser() io.ReadCloser {
	stringMessage := h.createHTTPMessage()
	fn := "message.http"
	fs := fstest.MapFS{fn: {Data: []byte(stringMessage)}}
	fh, _ := fs.Open(fn)
	defer fh.Close()
	return fh
}

func TestParsingHTTPGet(t *testing.T) {
	message := HTTPMessage{
		StartLine: "GET /coffee HTTP/1.1",
		FieldLines: []string{
			"Host: localhost:42069",
			"User-Agent: curl",
			"Accept: */*",
		},
		Body: "",
	}
	stream := message.httpMessageAsReadCloser()
	var got []string
	for line := range GetLinesChannel(stream) {
		got = append(got, line)
	}
	expected := message.expectedLines()
	result := slices.Compare(got, expected)
	if result != 0 {
		t.Errorf("Slices are different - got: %v | expected: %v\n", got, expected)
	}
}

func TestParsingHTTPPost(t *testing.T) {
	message := HTTPMessage{
		StartLine: "POST /coffee HTTP/1.1",
		FieldLines: []string{
			"Host: localhost:42069",
			"User-Agent: curl/8.16.0",
			"Accept: */*",
			"Content-Type: application/json",
			"Content-Length: 22",
		},
		Body: "{\"flavor\": \"dark mode\"}",
	}
	stream := message.httpMessageAsReadCloser()
	var got []string
	for line := range GetLinesChannel(stream) {
		got = append(got, line)
	}
	expected := message.expectedLines()
	result := slices.Compare(got, expected)
	if result != 0 {
		t.Errorf("Slices are different - got: %v | expected: %v\n", got, expected)
	}
}
