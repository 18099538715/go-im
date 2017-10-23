package server

import (
	"bean"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

var Udpchain = make(chan *bean.Protocol, 10000)

func init() {
	go func() {
		conn, err := net.Dial("udp", "127.0.0.1:9001")
		defer conn.Close()
		if err != nil {
			fmt.Println("udp客户端建立失败", err)
		}
		for {
			tmp, _ := <-Udpchain
			fmt.Println(tmp)
			b, _ := proto.Marshal(tmp)
			conn.Write(b)
		}

	}()
}
