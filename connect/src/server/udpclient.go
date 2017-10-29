package server

import (
	"fmt"
	"handle"
	"net"

	"github.com/golang/protobuf/proto"
)

func StartUdpClient() {
	go func() {
		conn, err := net.Dial("udp", "127.0.0.1:9001")
		defer conn.Close()
		if err != nil {
			fmt.Println("udp客户端建立失败", err)
		}
		for {
			tmp, _ := <-handle.Udpchain
			b, _ := proto.Marshal(tmp)
			conn.Write(b)
		}

	}()
}
