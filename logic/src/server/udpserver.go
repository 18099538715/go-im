package server

import (
	"bean"
	"config"
	"fmt"
	"handle"
	"net"
	"util"

	"github.com/golang/protobuf/proto"
)

type userInfo struct {
	UserId     int64
	DeviceType int32
	Onlineaddr string
}

func StartUdpServer() {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", util.GetLocalIP()+":"+config.GetUdpPort())
		if err != nil {
			fmt.Println("失败", err)
		}
		udpConn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("失败", err)
		}
		defer func() {
			udpConn.Close()
			fmt.Println("udp服务器失败")
		}()

		var buff = make([]byte, 100000)

		for {
			n, remoteAddr, err := udpConn.ReadFromUDP(buff)
			if err != nil {
				fmt.Println("", err)
			} else {
				buffer := buff[0:n]
				udpPkg := &bean.UdpProtPkg{}
				proto.Unmarshal(buffer, udpPkg)
				go handle.Handle(udpPkg, remoteAddr, udpConn)
			}

		}
	}()
}
