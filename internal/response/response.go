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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return nil
	case StatusBadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return nil
	case StatusInternalServerError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return nil
	default:
		w.Write([]byte(""))
		return nil
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		data := k + ": " + v + "\r\n"
		w.Write([]byte(data))
	}
	return nil
}
