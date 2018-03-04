package server

import (
	"bean"
	"fmt"
	"handle"
	"net"

	"github.com/golang/protobuf/proto"
)

type userInfo struct {
	UserId     int64
	DeviceType int32
	Onlineaddr string
}

func StartUdpServer() {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", ":9001")
		if err != nil {
			fmt.Println("失败", err)
		}
		udplisten, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("失败", err)
		}
		defer func() {
			udplisten.Close()
			fmt.Println("udp服务器失败")
		}()

		var buff = make([]byte, 10000)

		for {
			fmt.Println("逻辑层开始接收消息")
			n, remoteAddr, err := udplisten.ReadFromUDP(buff)
			if err != nil {
				fmt.Println("", err)
			} else {
				buffer := buff[0:n]
				udpPkg := &bean.UdpProtPkg{}
				proto.Unmarshal(buffer, udpPkg)
				handle.Handle(udpPkg, remoteAddr)
			}

		}
	}()
}
