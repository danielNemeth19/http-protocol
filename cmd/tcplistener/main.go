package main

import (
	"fmt"
	"github.com/danielNemeth19/http-protocol/internal/request"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic("Listener cannot be setup")
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic("listener cannot accept connections")
		}
		fmt.Println("Connection accepted")
		request, err := request.RequestFromReader(c)
		if err != nil {
			panic("Error parsing request line")
		}
		fmt.Printf("Request line:\n - Method: %s\n - Target: %s\n - Version: %s\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)
		fmt.Println("Headers:")
		for k, v := range request.Headers {
			fmt.Printf(" - %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(request.Body))
		fmt.Println("Connection closed")
	}
}
