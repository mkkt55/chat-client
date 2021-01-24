package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Header struct {
	flag    int32
	protoId int32
	bodyLen int32
}

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

func SendProto(m protoreflect.ProtoMessage, id ProtoId) bool {
	result, err := proto.Marshal(m)
	fmt.Print(result)
	if err != nil {
		fmt.Printf("Proto marshal error... %s\n", err.Error())
		return false
	}
	var bytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, 100)
	nn, err := writer.Write(bytes)
	err = writer.Flush()
	time.Sleep(time.Second * 3)
	binary.BigEndian.PutUint32(bytes, uint32(id))
	nn, err = writer.Write(bytes)
	err = writer.Flush()
	time.Sleep(time.Second * 3)
	binary.BigEndian.PutUint32(bytes, uint32(len(result)))
	nn, err = writer.Write(bytes)
	err = writer.Flush()
	time.Sleep(time.Second * 3)
	// var head []byte = make([]byte, 0)
	// head = append(head, len(result))
	// _, err = writer.Write(head)
	if err != nil {
		fmt.Printf("Write header error... %s\n", err.Error())
		return false
	}
	nn, err = writer.Write(result)
	fmt.Print("nn", nn)
	if err != nil {
		fmt.Printf("Write error... %s\n", err.Error())
		return false
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Flush error... %s\n", err.Error())
		return false
	}
	return true
}

func TestSend() bool {
	time.Sleep(time.Second * 3)
	var pack LoginReq
	str := "111"
	pack.Auth = &str
	SendProto(&pack, pack.GetId())
	str = "222"
	pack.Auth = &str
	SendProto(&pack, pack.GetId())

	time.Sleep(time.Second * 3)
	return true
}

func Read(buffer []byte) bool {
	_, err := reader.Read(buffer)
	if err != nil {
		fmt.Printf("Read error... %s\n", err.Error())
		return false
	}
	return true
}
