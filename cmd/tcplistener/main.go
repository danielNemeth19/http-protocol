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
	currEndLine := false
	isBody := false
	go func() {
		defer f.Close()
		for {
			n, err := f.Read(buf)
			if n > 0 {
				if isBody {
					part := string(buf[:n])
					curr_line += part
					continue
				}
				parts := strings.Split(string(buf[:n]), endLine)
				if len(parts) == 1 {
					curr_line += parts[0]
					currEndLine = false
				} else if len(parts) == 2 {
					if parts[0] == "" {
						if currEndLine {
							isBody = true
							ch <- endLine
							curr_line = parts[1]
							continue
						} else {
							currEndLine = true
						}
					}
					curr_line += parts[0]
					ch <- curr_line
					if parts[1] == "" {
						currEndLine = true
					}
					curr_line = parts[1]
				} else if len(parts) == 3 {
					curr_line += parts[0]
					ch <- curr_line
					ch <- endLine
					curr_line = parts[2]
				}
			}
			if err != nil {
				if isBody {
					ch <- curr_line
				}
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
