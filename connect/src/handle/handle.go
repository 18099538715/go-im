package handle

import (
	"bean"
	"fmt"
	"net"
	"rediscache"

	"github.com/golang/protobuf/proto"
)

var Udpchan = make(chan *bean.UdpProtPkg, 10000)

var UserCache map[int64]*net.TCPConn = make(map[int64]*net.TCPConn)
var UserCacheConn map[*net.TCPConn]int64 = make(map[*net.TCPConn]int64)

/**
**/
func Handle(tcpPkg *bean.TcpProtPkg, conn *net.TCPConn) {
	fmt.Println("接入层接收到消息", tcpPkg.GetPkgType(), len(tcpPkg.GetContent()))
	switch bean.PkgTypeEnum(tcpPkg.PkgType) {
	case bean.PkgTypeEnum_LOGIN_REQ:
		loginHandle(tcpPkg, conn)
	case bean.PkgTypeEnum_PING:
		beatHandle(tcpPkg, conn)
	default:
		transferHandle(tcpPkg, conn)
	}

}

/**
登录的处理
**/
func loginHandle(tcpPkg *bean.TcpProtPkg, conn *net.TCPConn) {
	loginMsg := &bean.LoginReq{}
	err := proto.Unmarshal(tcpPkg.GetContent(), loginMsg)
	if err != nil {
		fmt.Println("解码登录消息出错", err)
		conn.Close()
		return
	}

	u := &bean.UserInfo{UserId: loginMsg.GetUserId(), DeviceType: loginMsg.GetDeviceType(), OnlineIp: "127.0.0.1", Port: 9000}
	_, err = rediscache.SetOnlineUser(u)
	if err != nil {
		fmt.Println(err)
		CloseConn(conn)
		return
	}
	fmt.Println("用户登录:", loginMsg.GetUserId())
	UserCache[loginMsg.GetUserId()] = conn
	UserCacheConn[conn] = loginMsg.GetUserId()
}

/**
心跳的处理
**/
func beatHandle(tcpPkg *bean.TcpProtPkg, conn *net.TCPConn) {
	fmt.Println("心跳消息包")
}

/**
转发到逻辑层的包处理
**/
func transferHandle(tcpPkg *bean.TcpProtPkg, conn *net.TCPConn) {
	fmt.Println(len(tcpPkg.GetContent()))
	udpPkg := &bean.UdpProtPkg{PkgType: tcpPkg.GetPkgType(), Content: tcpPkg.GetContent()}
	Udpchan <- udpPkg
}
func CloseConn(conn *net.TCPConn) {
	fmt.Println("触发关闭")
	userId, ok := UserCacheConn[conn]
	delete(UserCacheConn, conn)
	if ok {
		delete(UserCache, userId)
	}
	rediscache.DelOnlineUser(userId)
	conn.Close()
}
