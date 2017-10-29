package server

import (
	"bean"
	"fmt"
	"handle"
	"net"
	"util"

	"github.com/golang/protobuf/proto"
)

func StartUdpServer() {
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
			n, _, err := udplisten.ReadFromUDP(buff)
			if err != nil {
				fmt.Println("udp接收包错误", err)
			}
			buffer := buff[0:n]
			protocol := &bean.UdpProtocol{}
			proto.Unmarshal(buffer, protocol)

			if conn, ok := handle.UserCache[protocol.GetToUuserId()]; ok {
				conn.Write(util.Unit16Tobyte(uint16(protocol.GetProtocolType())))
				var l uint32 = uint32(len(protocol.GetProtocolContent()))

				conn.Write(util.Unit32Tobyte(l))          //写入长度
				conn.Write(protocol.GetProtocolContent()) //写入内容
			} else {
				fmt.Println("通道不存在")
			}

		}
	}()
}
