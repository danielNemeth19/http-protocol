package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fh, err := os.Open("messages.txt")
	if err != nil {
		panic("Problem opening the file")
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	buf := make([]byte, 4096)
	scanner.Buffer(buf, 8)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) == 8 {
			return 8, data[:8], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}
