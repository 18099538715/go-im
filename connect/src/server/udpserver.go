package server

import (
	"bean"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

var UserCache map[int64]*net.TCPConn

func init() {
	fmt.Println("udpserver初始化")
	UserCache = make(map[int64]*net.TCPConn)
	go func() {
		addr, err := net.ResolveUDPAddr("udp", ":9000")
		if err != nil {
			fmt.Println("", err)
		}
		udplisten, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("", err)
		}
		defer udplisten.Close()
		var buff = make([]byte, 10240)
		for {
			fmt.Println("开始读取")
			n, _, err := udplisten.ReadFromUDP(buff)
			if err != nil {
				fmt.Println("udp接收包错误", err)
			}
			buffer := buff[0:n]
			protocol := &bean.Protocol{}
			fmt.Println("来自逻辑层的消息", protocol)
			proto.Unmarshal(buffer, protocol)

			if conn, ok := UserCache[protocol.GetUserId()]; ok {
				conn.Write(protocol.GetProtocolContent())
			} else {
				fmt.Println("通道不存在")
			}

		}
	}()
}
