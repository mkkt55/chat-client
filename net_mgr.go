package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoPack struct {
	flag    byte
	protoId uint32
	bodyLen uint32
	body    []byte
}

const (
	headerLen = 9
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

func SendProto(m protoreflect.ProtoMessage, id ProtoId) error {
	result, err := proto.Marshal(m)
	if err != nil {
		fmt.Printf("Proto marshal error... %s\n", err.Error())
		return err
	}
	var flag byte
	header, err := buildHeader(flag, id, len(result))
	if err != nil {
		fmt.Printf("Build header error... %s\n", err.Error())
		return err
	}

	_, err = writer.Write(header)
	if err != nil {
		fmt.Printf("Write header error... %s\n", err.Error())
		return err
	}
	_, err = writer.Write(result)
	if err != nil {
		fmt.Printf("Write body error... %s\n", err.Error())
		return err
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Flush error... %s\n", err.Error())
		return err
	}
	return nil
}

func ReadProto() (*ProtoPack, error) {
	var pack ProtoPack
	bytes := make([]byte, headerLen)
	_, err := io.ReadFull(reader, bytes)
	if err != nil {
		return nil, err
	}
	fmt.Println(bytes)
	pack.flag = bytes[0]
	pack.protoId = binary.BigEndian.Uint32(bytes[1:5])
	pack.bodyLen = binary.BigEndian.Uint32(bytes[5:9])
	pack.body = make([]byte, pack.bodyLen)
	fmt.Println(pack)
	io.ReadFull(reader, pack.body)
	return &pack, nil
}

func buildHeader(flag byte, id ProtoId, bodyLen int) ([]byte, error) {
	var bytes []byte = make([]byte, 0)
	bytes = append(bytes, flag)
	binary.BigEndian.PutUint32(bytes, uint32(id))
	binary.BigEndian.PutUint32(bytes, uint32(bodyLen))
	return bytes, nil
}
