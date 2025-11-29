package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/danielNemeth19/http-protocol/internal/request"
	"github.com/danielNemeth19/http-protocol/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

func (h HandlerError) WriteError(w io.Writer) {
	response.WriteStatusLine(w, response.StatusBadRequest)
	headers := response.GetDefaultHeaders(0)
	response.WriteHeaders(w, headers)
	w.Write([]byte("\r\n"))
}

type Server struct {
	listener   net.Listener
	inShutdown atomic.Bool
}

func (s *Server) Close() error {
	s.inShutdown.Store(true)
	s.listener.Close()
	return nil
}

func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			isShutdown := s.inShutdown.Load()
			if isShutdown {
				break
			}
			log.Printf("Error during accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		errH := HandlerError{Message: err.Error(), Code: response.StatusBadRequest}
		errH.WriteError(conn)
	}
	var buf bytes.Buffer
	handlerError := handler(&buf, req)
	if handlerError != nil {
		response.WriteStatusLine(&buf, handlerError.Code)
		headers := response.GetDefaultHeaders(len(handlerError.Message))
		response.WriteHeaders(&buf, headers)
		buf.Write([]byte("\r\n"))
		buf.Write([]byte(handlerError.Message))
		buf.WriteTo(conn)
		return
	}
	var buf2 bytes.Buffer
	response.WriteStatusLine(&buf2, response.StatusOK)
	headers := response.GetDefaultHeaders(len(buf.String()))
	response.WriteHeaders(&buf2, headers)
	buf2.Write([]byte("\r\n"))
	buf.WriteTo(&buf2)
	buf2.WriteTo(conn)
}

func Serve(port int, handler Handler) (*Server, error) {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener}
	go server.listen(handler)
	return server, nil
}
