package server

import (
	"bean"
	"bytes"
	"conf"
	"encoding/binary"
	"fmt"
	"handle"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
)

func StartTcpServer() {
	go func() {
		//addr, err := net.ResolveTCPAddr("tcp", "192.168.31.248:8888")
		tcpServerPort := conf.ConfMap["tcpserverport"]

		addr, err := net.ResolveTCPAddr("tcp", ":"+tcpServerPort)
		if err != nil {
			fmt.Println("tcp启动失败", err)
			return
		}
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
	defer handle.CloseConn(conn)
	for {
		lengthByte := make([]byte, 4)
		var length int32
		_, err := io.ReadFull(conn, lengthByte)
		if err != nil {
			fmt.Println("读取信息出错", err)
			break
		}
		bytesBuffer := bytes.NewBuffer(lengthByte)
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		msgByte := make([]byte, length)
		conn.Read(msgByte)
		tcpPkg := &bean.TcpProtPkg{}
		err = proto.Unmarshal(msgByte, tcpPkg)
		if err != nil {
			fmt.Println("读取信息出错")
			break
		}
		handle.Handle(tcpPkg, conn)
	}

}
