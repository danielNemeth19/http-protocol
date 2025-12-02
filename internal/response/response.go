package response

import (
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

var BadRequest = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

var InternalServerError = `<html>
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

type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		w.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return nil
	case StatusBadRequest:
		w.Writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return nil
	case StatusInternalServerError:
		w.Writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return nil
	default:
		w.Writer.Write([]byte(""))
		return nil
	}
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
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
	headers.Set("Content-Type", "text/html")
	return headers
}
