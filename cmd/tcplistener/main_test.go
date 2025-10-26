package main

import (
	"io"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/danielNemeth19/http-protocol/internal/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HTTPMessage struct {
	StartLine  string
	FieldLines []string
	Body       string
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
	request, err := request.RequestFromReader(stream)
	require.NoError(t, err)
	require.NotNil(t, request)
	assert.Equal(t, request.RequestLine.Method, "GET")
	assert.Equal(t, request.RequestLine.RequestTarget, "/coffee")
	assert.Equal(t, request.RequestLine.HttpVersion, "1.1")
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
	request, err := request.RequestFromReader(stream)
	require.NoError(t, err)
	require.NotNil(t, request)
	assert.Equal(t, request.RequestLine.Method, "POST")
	assert.Equal(t, request.RequestLine.RequestTarget, "/coffee")
	assert.Equal(t, request.RequestLine.HttpVersion, "1.1")
}
