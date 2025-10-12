package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func GetLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	buf := make([]byte, 8)
	curr_line := ""
	endLine := "\r\n"
	// lastChar := ""
	go func() {
		defer f.Close()
		for {
			n, err := f.Read(buf)
			if n > 0 {
				// TODO: every "" needs to be ignored excepts if it a double one
				parts := strings.Split(string(buf[:n]), endLine)
				if len(parts) == 1 {
					curr_line += parts[0]
				} else if len(parts) == 2 {
					curr_line += parts[0]
					ch <- curr_line
					curr_line = parts[1]
				} else if len(parts) == 3 {
					curr_line += parts[0]
					ch <- curr_line
					ch <- "\r\n"
					curr_line = parts[2]
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
		for c := range GetLinesChannel(c) {
			fmt.Printf("%s\r\n", strings.TrimRight(c, "\r\n"))
		}
		fmt.Println("Connection closed")
	}
}
