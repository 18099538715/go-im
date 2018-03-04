package handle

import (
	"bean"
	"fmt"
	"net"
	"rediscache"

	"github.com/golang/protobuf/proto"
)

func Handle(udpPkg *bean.UdpProtPkg, reomteUdpAddr *net.UDPAddr) {
	switch bean.PkgTypeEnum(udpPkg.PkgType) {
	case bean.PkgTypeEnum_MSG_SEND_SINGLE: //单聊消息发送包
		singleMsg(udpPkg, reomteUdpAddr)
	case bean.PkgTypeEnum_MSG_SEND_GROUP: //群聊消息发送包
		groupMsg(udpPkg, reomteUdpAddr)
	case bean.PkgTypeEnum_MSG_REQ: //消息请求包
		msgReq(udpPkg, reomteUdpAddr)
	case bean.PkgTypeEnum_MSG_REQ_ACK: //消息请求ack包
		msgReqAck(udpPkg, reomteUdpAddr)
	default:
	}
}

/**
消息请求ack，用户拉取完消息给个ack，更新用户的syncKey
**/
func msgReqAck(udpPkg *bean.UdpProtPkg, reomteUdpAddr *net.UDPAddr) {
	msg := &bean.MsgReqAck{}
	err := proto.Unmarshal(udpPkg.GetContent(), msg)
	if err != nil {
		fmt.Println("reqAck解析包出错", err)
	}
	err = rediscache.SetCurrentUserSrlNo(msg.GetUserId(), msg.GetSrlNo())
	if err != nil {
		fmt.Println("更新用户syncKey出错")
	}
}

/**
处理消息请求
**/
func msgReq(udpPkg *bean.UdpProtPkg, reomteUdpAddr *net.UDPAddr) {
	msg := &bean.MsgReq{}
	err := proto.Unmarshal(udpPkg.GetContent(), msg)
	if err != nil {
		fmt.Println("拉取用户消息出错1", err)
		return
	}
	srlNos := make([]int64, msg.GetPageSize())
	currentSrlNo, err := rediscache.GetCurrentUserSrlNo(msg.GetUserId())
	if currentSrlNo < msg.GetSrlNo() {
		currentSrlNo = msg.GetSrlNo()
	}
	fmt.Println("当前用户的srlNo:", currentSrlNo)
	for index, _ := range srlNos {
		srlNos[index] = currentSrlNo + int64(index+1)
	}
	msgs, err := rediscache.GetUserMsgs(msg.GetUserId(), srlNos)
	if err != nil {
		fmt.Println("拉取用户消息出错3", err)
		return
	}
	msgReqRes := &bean.MsgReqRes{}
	if msgs != nil {
		msgReqRes.Msgs = msgs
	}
	msgReqRes.SrlNo = currentSrlNo + int64(len(srlNos))
	content, err := proto.Marshal(msgReqRes)
	res := &bean.UdpProtPkg{ToUserId: msg.GetUserId(), PkgType: int32(bean.PkgTypeEnum_MSG_REQ_RES)}
	if content != nil {
		res.Content = content
	}
	b, err := proto.Marshal(res)
	if err != nil {
		fmt.Println("拉取用户消息出错5", err)
		return
	}
	udpclient, err := net.DialUDP("udp", nil, reomteUdpAddr)
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	defer udpclient.Close()
	udpclient.Write(b)
	udpclient.Close()
}

/**
处理群组消息
**/
func groupMsg(udpPkg *bean.UdpProtPkg, reomteUdpAddr *net.UDPAddr) {

}

/**
处理单聊消息
**/
func singleMsg(udpPkg *bean.UdpProtPkg, reomteUdpAddr *net.UDPAddr) {
	singleMsg := &bean.SingleMsg{}
	proto.Unmarshal(udpPkg.GetContent(), singleMsg)
	fmt.Println("传过来的消息", singleMsg)
	srlNo, err := rediscache.IncrUserSrlNo(singleMsg.GetToUserId())
	if err != nil {
		fmt.Println("srlNo自增失败", err)
		msgSendRes(singleMsg.GetFromUserId(), reomteUdpAddr, singleMsg.GetMsgId(), false)
		return
	}
	singleMsg.SrlNo = srlNo
	err = rediscache.SetUserSingleMsg(singleMsg)
	if err != nil {
		fmt.Println("保存消息失败", err)
		msgSendRes(singleMsg.GetFromUserId(), reomteUdpAddr, singleMsg.GetMsgId(), false)
		return
	}
	//给发送者发响应
	msgSendRes(singleMsg.GetFromUserId(), reomteUdpAddr, singleMsg.GetMsgId(), true)
	//给接收者发通知
	msgInform(singleMsg.GetToUserId(), srlNo)
}

/**
消息发送响应
**/
func msgSendRes(userId int64, reomteUdpAddr *net.UDPAddr, msgId string, flag bool) {
	udpclient, err := net.DialUDP("udp", nil, reomteUdpAddr)
	defer udpclient.Close()
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	msgSendRes := &bean.MsgSendRes{Flag: true, MsgId: msgId}
	content, _ := proto.Marshal(msgSendRes)
	res := &bean.UdpProtPkg{ToUserId: userId, Content: content, PkgType: int32(bean.PkgTypeEnum_MSG_SEND_RES)}
	b, _ := proto.Marshal(res)
	udpclient.Write(b)
	udpclient.Close()
}

/**
新消息通知
**/
func msgInform(userId int64, srlNo int64) {
	u := rediscache.GetOnlineUser(userId)
	if u == nil {
		fmt.Println("用户不在线", userId)
		return
	}
	udpclient, err := net.Dial("udp", u.Onlineaddr)
	defer udpclient.Close()
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	inForm := &bean.MsgInform{SrlNo: srlNo, UserId: userId}
	b, _ := proto.Marshal(inForm)
	res := &bean.UdpProtPkg{ToUserId: userId, Content: b, PkgType: int32(bean.PkgTypeEnum_MSG_INFORM)}
	t, _ := proto.Marshal(res)
	udpclient.Write(t)
	udpclient.Close()
}
