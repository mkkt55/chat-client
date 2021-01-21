package main

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var conn net.Conn
var reader *bufio.Reader
var writer *bufio.Writer

func InitConnection() bool {
	var err error
	conn, err = net.Dial("tcp4", "192.168.93.233:15000")
	if err != nil {
		fmt.Printf("Connect error...\n")
		return false
	}
	reader = bufio.NewReader(conn)
	writer = bufio.NewWriter(conn)
	return true
}

func ReleaseConnection() bool {
	err := conn.Close()
	if err != nil {
		fmt.Printf("Close connection error...\n")
		return false
	}
	return true
}

func SendProto(m protoreflect.ProtoMessage) bool {
	result, err := proto.Marshal(m)
	if err != nil {
		fmt.Printf("Proto marshal error... %s\n", err.Error())
		return false
	}
	// var head []byte = make([]byte, 0)
	// head = append(head, len(result))
	// _, err = writer.Write(head)
	if err != nil {
		fmt.Printf("Write error... %s\n", err.Error())
		return false
	}
	_, err = writer.Write(result)
	if err != nil {
		fmt.Printf("Write error... %s\n", err.Error())
		return false
	}
	return true
}

func TestSend() bool {
	time.Sleep(time.Second * 5)
	var pack LoginReq
	fmt.Print(pack.GetId(), "\n")
	return SendProto(&pack)
}

func Read(buffer []byte) bool {
	_, err := reader.Read(buffer)
	if err != nil {
		fmt.Printf("Read error... %s\n", err.Error())
		return false
	}
	return true
}
