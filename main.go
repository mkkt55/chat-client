package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp4", "192.168.93.233:15000")
	if err != nil {
		fmt.Printf("Connect error...\n")
		return
	}
	status, err := bufio.NewReader(conn).ReadString(' ')
	fmt.Print(status)
	fmt.Fprintf(conn, "Hi")
}
