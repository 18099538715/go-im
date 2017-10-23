package handle

import (
	"bean"
	"encoding/json"
	"fmt"
	"net"
	"redis"
	"server"

	"github.com/golang/protobuf/proto"
)

type userInfo struct {
	UserId     int64
	DeviceType int32
	Onlineaddr string
}

/**
**/
func Handle(msg *bean.Protocol, conn *net.TCPConn) {
	switch bean.ProtocolTypeEnum(msg.ProtocolType) {
	case bean.ProtocolTypeEnum_LOGIN:
		loginHandle(msg, conn)
	case bean.ProtocolTypeEnum_PONG:
		beatHandle(msg, conn)
	default:
		transferHandle(msg, conn)
	}

}

/**
登录的处理
**/
func loginHandle(msg *bean.Protocol, conn *net.TCPConn) {
	loginMsg := &bean.LoginReq{}
	err := proto.Unmarshal(msg.GetProtocolContent(), loginMsg)
	if err != nil {
		fmt.Println("解码消息出错", err)
		return
	}
	u := &userInfo{UserId: loginMsg.GetUerId(), DeviceType: loginMsg.GetDeviceType(), Onlineaddr: "127.0.0.1:9000"}
	fmt.Println(u)
	a, _ := json.Marshal(u)
	fmt.Println(string(a))
	_, err = redis.Cluster.Do("HSET", "onlineinfo", loginMsg.GetUerId(), string(a))
	if err != nil {
		fmt.Println("登录出错", err)
		conn.Close()
		return
	}
	server.UserCache[loginMsg.GetUerId()] = conn

}

/**
心跳的处理
**/
func beatHandle(msg *bean.Protocol, conn *net.TCPConn) {
	fmt.Println("心跳消息包")
}

/**
转发到逻辑层的包处理
**/
func transferHandle(msg *bean.Protocol, conn *net.TCPConn) {
	server.Udpchain <- msg
}
