package handle

import (
	"bean"
	"fmt"
	"net"
	"redis"

	"github.com/golang/protobuf/proto"
)

var Udpchain = make(chan *bean.UdpProtocol, 10000)

var UserCache map[int64]*net.TCPConn = make(map[int64]*net.TCPConn)

/**
**/
func Handle(msg *bean.TcpProtocol, conn *net.TCPConn) {
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
func loginHandle(msg *bean.TcpProtocol, conn *net.TCPConn) {

	loginMsg := &bean.LoginReq{}
	err := proto.Unmarshal(msg.GetProtocolContent(), loginMsg)
	if err != nil {
		fmt.Println("解码消息出错", err)
		return
	}
	u := &bean.UserInfo{UserId: loginMsg.GetUerId(), DeviceType: loginMsg.GetDeviceType(), Onlineaddr: "127.0.0.1:9000"}
	_, err = redis.SetOnlineUser(u)
	if err != nil {
		fmt.Println("登录出错", err)
		conn.Close()
		return
	}
	fmt.Println("用户登录:", loginMsg.GetUerId())
	UserCache[loginMsg.GetUerId()] = conn
}

/**
心跳的处理
**/
func beatHandle(msg *bean.TcpProtocol, conn *net.TCPConn) {
	fmt.Println("心跳消息包")
}

/**
转发到逻辑层的包处理
**/
func transferHandle(msg *bean.TcpProtocol, conn *net.TCPConn) {
	udpMsg := &bean.UdpProtocol{ProtocolType: msg.GetProtocolType(), ProtocolContent: msg.GetProtocolContent(), Fromaddress: "127.0.0.1:9000"}
	Udpchain <- udpMsg
}
