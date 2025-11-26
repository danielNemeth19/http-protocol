package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/danielNemeth19/http-protocol/internal/response"
)

type Server struct {
	listener   net.Listener
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
	// response := []byte(
	// "HTTP/1.1 200 OK\r\n" +
	// "Content-Type: text/plain\r\n" +
	// "Content-Length: 13\r\n\r\n" +
	// "Hello World!\n",
	// )
	// conn.Write(response)
	response.WriteStatusLine(conn, 200)
	headers := response.GetDefaultHeaders(13)
	response.WriteHeaders(conn, headers)
}

func Serve(port int) (*Server, error) {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener}
	go server.listen()
	return server, nil
}
