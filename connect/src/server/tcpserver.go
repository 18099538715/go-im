package server

import (
	"bean"
	"bytes"
	"encoding/binary"
	"fmt"
	"handle"
	"io"
	"net"
	"time"
)

func StartTcpServer() {
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", "192.168.31.248:8888")
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
			go HandleConnRead(conn)
		}
	}()

}

func HandleConnRead(conn *net.TCPConn) {
	defer func() {
		conn.Close()
	}()

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
		bytesBuffer := bytes.NewBuffer(protocolTypeByte)
		binary.Read(bytesBuffer, binary.BigEndian, &protocolType)
		io.ReadFull(conn, lengthByte)
		bytesBuffer = bytes.NewBuffer(lengthByte)
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		msgByte := make([]byte, length)
		conn.Read(msgByte)
		protocol := &bean.TcpProtocol{ProtocolContent: msgByte, ProtocolType: int32(protocolType)}
		handle.Handle(protocol, conn)
	}

}
