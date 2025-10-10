package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		panic("UPD address setup failed")
	}
	connection, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic("UDP connection failed")
	}
	defer connection.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">")
		line, err := reader.ReadString(byte('\n'))
		if err != nil {
			panic("cannot read line")
		}
		connection.Write([]byte(line))
	}
}
