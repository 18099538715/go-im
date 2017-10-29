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
				fmt.Println("", err)
			} else {
				buffer := buff[0:n]
				protocol := &bean.UdpProtocol{}
				proto.Unmarshal(buffer, protocol)
				handle.Handle(protocol)
			}

		}
	}()
}
