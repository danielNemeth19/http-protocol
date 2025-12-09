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
	done
)

type Writer struct {
	Writer io.Writer
	state  writeState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != initalized {
		return fmt.Errorf("Writer expected to be initialized, got: %d", w.state)
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
	if w.state != writeStateHeaders {
		return fmt.Errorf("Writer expected to be in writeStateHeaders, got: %d", w.state)
	}
	for k, v := range headers {
		data := k + ": " + v + "\r\n"
		w.Writer.Write([]byte(data))
	}
	w.Writer.Write([]byte("\r\n"))
	w.state = writeStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writeStateBody {
		return 0, fmt.Errorf("Writer expected to be in writeStateBody state, got: %d", w.state)
	}
	w.Writer.Write(p)
	w.state = done
	return len(p), nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	hexa := fmt.Sprintf("%X", len(p))
	w.Writer.Write([]byte(hexa))
	w.Writer.Write([]byte("\r\n"))
	w.Writer.Write(p)
	w.Writer.Write([]byte("\r\n"))
	return 0, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	w.Writer.Write([]byte(fmt.Sprintf("%X", 0)))
	w.Writer.Write([]byte("\r\n"))
	w.Writer.Write([]byte("\r\n"))
	return 0, nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func GetChunkedHeaders() headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Type", "text/plain")
	headers.Set("Transfer-Encoding", "chunked")
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
