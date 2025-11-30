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
	response.WriteStatusLine(w, h.Code)
	headers := response.GetDefaultHeaders(len(h.Message))
	response.WriteHeaders(w, headers)
	w.Write([]byte("\r\n"))
	w.Write([]byte(h.Message))
}

type Server struct {
	listener   net.Listener
	handler    Handler
	inShutdown atomic.Bool
}

func (s *Server) Close() error {
	s.inShutdown.Store(true)
	s.listener.Close()
	return nil
}

func (s *Server) listen() {
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
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		errH := HandlerError{Message: err.Error(), Code: response.StatusBadRequest}
		errH.WriteError(conn)
		return
	}
	var body bytes.Buffer
	handlerError := s.handler(&body, req)
	if handlerError != nil {
		handlerError.WriteError(conn)
		return
	}
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(body.Len())
	response.WriteHeaders(conn, headers)
	conn.Write([]byte("\r\n"))
	body.WriteTo(conn)
}

func Serve(port int, handler Handler) (*Server, error) {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	go server.listen()
	return server, nil
}
