package main

import (
	"fmt"
	"io"
	"os"
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
	fh, err := os.Open("messages.txt")
	if err != nil {
		panic("Problem opening the file")
	}
	for c := range getLinesChannel(fh) {
		fmt.Printf("read: %s\n", c)
	}
}
