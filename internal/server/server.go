package server

import (
	"net"
	"strconv"
)

type Server struct {
	state bool
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() error {
	return nil
}

func (s *Server) handle(conn net.Conn) {
	response := []byte(
		"HTTP/1.1 200 OK\r\n" +
			"Content-Type: test/plain\r\n" +
			"Content-Length: 13\r\n\r\n" +
			"Hello World!",
	)
	conn.Write(response)
}

func Serve(port int) (*Server, error) {
	server := Server{state: true}
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &server, nil
}
