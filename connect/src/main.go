package main

import (
	"bean"
	"bytes"
	"encoding/binary"
	"fmt"
	"handle"
	"io"
	"net"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr, err := net.ResolveTCPAddr("tcp", ":8888")
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen error: ", addr)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
			break
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(10 * time.Minute)
		// start a new goroutine to handle the new connection
		go HandleConnRead(conn)
	}

}
func HandleConnRead(conn *net.TCPConn) {
	defer conn.Close()

	for {
		protocolTypeByte := make([]byte, 2)
		lengthByte := make([]byte, 4)
		var protocolType int16
		var length int32
		_, err := io.ReadFull(conn, protocolTypeByte)
		if err != nil {
			conn.Close()
			fmt.Println("通道出错", err)
			break
		}
		fmt.Println(protocolTypeByte)
		bytesBuffer := bytes.NewBuffer(protocolTypeByte)
		binary.Read(bytesBuffer, binary.BigEndian, &protocolType)
		fmt.Println("消息报类型", protocolType)
		io.ReadFull(conn, lengthByte)
		bytesBuffer = bytes.NewBuffer(lengthByte)
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		fmt.Println("消息长度", length)
		msgByte := make([]byte, length)
		conn.Read(msgByte)
		/**msg := &bean.SingleMsg{}
		err := proto.Unmarshal(msgByte, msg)
		if err != nil {
			fmt.Println("解码错误", err)
			conn.Close()
			break
		}
		fmt.Println("接收到的消息", msg)**/
		protocol := &bean.Protocol{ProtocolContent: msgByte, ProtocolType: int32(protocolType)}
		handle.Handle(protocol, conn)
	}

}
