package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	buf := make([]byte, 8)
	curr_line := ""
	go func() {
		defer f.Close()
		for {
			n, err := f.Read(buf)
			if n > 0 {
				parts := strings.Split(string(buf[:n]), "\n")
				if len(parts) == 1 {
					curr_line += parts[0]
				} else if len(parts) == 2 {
					curr_line += parts[0]
					ch <- curr_line
					curr_line = parts[1]
				}
			}
			if err != nil {
				break
			}
		}
		close(ch)
	}()
	return ch
}

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
		for c := range getLinesChannel(c) {
			fmt.Printf("read: %s\n", c)
		}
		fmt.Println("Connection closed")
	}
}
