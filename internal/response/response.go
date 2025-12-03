package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/danielNemeth19/http-protocol/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var BadRequestHTML = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

var InternalServerErrorHTML = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

var SuccessHTML = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

type writeState int

const (
	initalized writeState = iota
	writeStateHeaders
	writeStateBody
)

type Writer struct {
	Writer io.Writer
	state  writeState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != initalized {
		return fmt.Errorf("Writer expected to be initialized, got: %s", w.state)
	}
	var statusLine string
	switch statusCode {
	case StatusOK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusBadRequest:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case StatusInternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		return fmt.Errorf("Unrecognized status code: %d\n", statusCode)
	}
	_, err := w.Writer.Write([]byte(statusLine))
	if err != nil {
		return err
	}
	w.state = writeStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writeStateStatusLine {
		return fmt.Errorf("Writer expected to be initialized, got: %s", w.state)
	}
	w.state = writeStateStatusLine
	for k, v := range headers {
		data := k + ": " + v + "\r\n"
		w.Writer.Write([]byte(data))
	}
	w.Writer.Write([]byte("\r\n"))
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.Writer.Write(p)
	return len(p), nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func ReplaceHeader(header, headers headers.Headers) headers.Headers {
	for key, value := range header {
		if _, exists := headers[key]; exists {
			headers[key] = value
		}
	}
	return headers
}
