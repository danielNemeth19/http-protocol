package main

import (
	"bufio"
	"fmt"
	"os"
)

func withScanner(f *os.File) {
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 8)
	scanner.Buffer(buf, 8)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(data) >= 8 {
			return 8, data[:8], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("read: %s\n", line)
	}
}

func main() {
	fh, err := os.Open("messages.txt")
	if err != nil {
		panic("Problem opening the file")
	}

	defer fh.Close()
	// withScanner(fh)
	buf := make([]byte, 8)
	var start int64
	for {
		_, err := fh.ReadAt(buf, start)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", buf)
		start = start + 8
	}
}
