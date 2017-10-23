package server

import (
	"bean"
	"encoding/json"
	"fmt"
	"net"
	"redis"

	"github.com/golang/protobuf/proto"
)

type userInfo struct {
	UserId     int64
	DeviceType int32
	Onlineaddr string
}

var Tmp int

func init() {
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
		}
		buffer := buff[0:n]
		fmt.Println("收到消息的长度**", n)
		protocol := &bean.Protocol{}
		proto.Unmarshal(buffer, protocol)
		msg := &bean.SingleMsg{}
		proto.Unmarshal(protocol.GetProtocolContent(), msg)
		reply, err := redis.Cluster.Do("HGET", "onlineinfo", msg.GetFromUserId())
		value, ok := reply.([]byte)
		if !ok {
			fmt.Println("类型错误")
		}
		u1 := &userInfo{}
		err = json.Unmarshal(value, u1)
		fmt.Println("redis读取", u1)

		udpclient, err := net.Dial("udp", u1.Onlineaddr)
		defer udpclient.Close()
		if err != nil {
			fmt.Println("udp客户端建立失败", err)
		}
		toMsg := &bean.MsgSendRes{Flag: true, MsgId: "1"}
		b, _ := proto.Marshal(toMsg)
		protocolRes := &bean.Protocol{UserId: msg.GetFromUserId(), ProtocolContent: b}
		t, _ := proto.Marshal(protocolRes)
		fmt.Println("像接入层输出消息")
		udpclient.Write(t)
		udpclient.Close()
	}
}
