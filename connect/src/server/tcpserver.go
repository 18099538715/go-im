package server

import (
	"bean"
	"bytes"
	"config"
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

		addr, err := net.ResolveTCPAddr("tcp", config.GetTcpIp()+":"+config.GetTcpPort())
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
		connTimeout := config.GetConnTimeOut()
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("accept error: ", err)
				break
			}
			go HandleConnRead(conn, connTimeout)
		}
	}()

}

func HandleConnRead(conn *net.TCPConn, connTimeOut int64) {
	defer handle.CloseConn(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute * time.Duration(connTimeOut)))
		lengthByte := make([]byte, 4)
		var length int32
		_, err := io.ReadFull(conn, lengthByte)
		if err != nil {
			fmt.Println("读取信息出错1", err)
			break
		}
		bytesBuffer := bytes.NewBuffer(lengthByte)
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		msgByte := make([]byte, length)
		_, err = conn.Read(msgByte)
		if err != nil {
			fmt.Println("读取信息出错2", err)
			break
		}
		tcpPkg := &bean.TcpProtPkg{}
		err = proto.Unmarshal(msgByte, tcpPkg)
		if err != nil {
			fmt.Println("读取信息出错3")
			break
		}
		handle.Handle(tcpPkg, conn)
	}

}
