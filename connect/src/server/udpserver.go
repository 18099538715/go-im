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

func StartUdpServer() {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", util.GetLocalIP()+":"+config.GetUdpPort())
		if err != nil {
			fmt.Println("建立连接失败", err)
		}
		udplisten, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Println("", err)
		}
		go receive(udplisten)
		go send(udplisten)

	}()
}
func receive(udpconn *net.UDPConn) {
	var buff = make([]byte, 10240)
	for {
		n, _, err := udpconn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println("udp接收包错误", err)
			continue
		}
		buffer := buff[0:n]
		pkg := &bean.UdpProtPkg{}
		err = proto.Unmarshal(buffer, pkg)
		if err != nil {
			fmt.Println("udp接收包错误", err)
			continue
		}
		fmt.Println("逻辑层发过来的消息：", pkg.PkgType)
		tcpPkg := &bean.TcpProtPkg{PkgType: pkg.PkgType, Content: pkg.Content}
		if conn, ok := handle.UserCache[pkg.GetToUserId()]; ok {
			b, _ := proto.Marshal(tcpPkg)
			var l uint32 = uint32(len(b))
			conn.Write(util.Unit32Tobyte(l)) //写入长度
			conn.Write(b)                    //写入内容
		} else {
			fmt.Println("通道不存在")
		}

	}
}
func send(conn *net.UDPConn) {
	ip := net.ParseIP(config.GetRemoteUdpIp())
	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: int(config.GetRemoteUdpPort()),
	}
	for {
		tmp, ok := <-handle.Udpchan
		if !ok {
			fmt.Println("udpchain获取消息失败")
			continue
		}
		b, _ := proto.Marshal(tmp)
		conn.WriteToUDP(b, udpAddr)
		fmt.Println("往逻辑层发送消息")
	}
}
