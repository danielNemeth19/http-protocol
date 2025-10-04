package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fh, err := os.Open("messages.txt")
	if err != nil {
		panic("Problem opening the file")
	}
	defer fh.Close()

	buf := make([]byte, 8)
	var start int64
	curr_line := ""
	for {
		n, err := fh.ReadAt(buf, start)
		if n > 0 {
			parts := strings.Split(string(buf[:n]), "\n")
			if len(parts) == 1 {
				curr_line += parts[0]
			} else if len(parts) == 2 {
				curr_line += parts[0]
				fmt.Printf("read: %s\n", curr_line)
				curr_line = parts[1]
			}
			start += int64(n)
		}
		if err != nil {
			break
		}
	}
}
