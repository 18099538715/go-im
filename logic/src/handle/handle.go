package handle

import (
	"bean"
	"fmt"
	"net"
	"redis"

	"github.com/golang/protobuf/proto"
)

func Handle(msg *bean.UdpProtocol) {
	switch bean.ProtocolTypeEnum(msg.ProtocolType) {
	case bean.ProtocolTypeEnum_MSG_SEND_SINGLE: //单聊消息发送包
		singleMsg(msg)
	case bean.ProtocolTypeEnum_MSG_SEND_GROUP: //群聊消息发送包
		groupMsg(msg)
	case bean.ProtocolTypeEnum_MSG_REQ: //消息请求包
		msgReq(msg)
	case bean.ProtocolTypeEnum_MSG_REQ_ACK: //消息请求ack包
		msgReqAck(msg)
	default:

	}
}

/**
消息请求ack，用户拉取万消息给个ack，更新用户的syncKey
**/
func msgReqAck(protocol *bean.UdpProtocol) {
	msg := &bean.MsgReqAck{}
	err := proto.Unmarshal(protocol.GetProtocolContent(), msg)
	if err != nil {
		fmt.Println("reqAck解析包出错", err)
	}
	err = redis.SetCurrentUserSynckey(msg.GetUserId(), msg.GetSyncKey())
	if err != nil {
		fmt.Println("更新用户syncKey出错")
	}
}

/**
处理消息请求
**/
func msgReq(protocol *bean.UdpProtocol) {
	msg := &bean.MsgReq{}
	err := proto.Unmarshal(protocol.GetProtocolContent(), msg)
	if err != nil {
		fmt.Println("拉取用户消息出错", err)
		return
	}
	synckeys := make([]int64, msg.GetPageSize())
	redis.GetUserMsgs(msg.GetUserId(), synckeys)
	udpclient, err := net.Dial("udp", protocol.GetFromaddress())
	defer udpclient.Close()
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	msgReqRes := &bean.MsgReqRes{}
	content, _ := proto.Marshal(msgReqRes)
	protocolRes := &bean.UdpProtocol{ToUuserId: msg.GetUserId(), ProtocolContent: content, ProtocolType: int32(bean.ProtocolTypeEnum_MSG_REQ_RES)}
	b, _ := proto.Marshal(protocolRes)
	udpclient.Write(b)
	udpclient.Close()
}

/**
处理群组消息
**/
func groupMsg(protocol *bean.UdpProtocol) {

}

/**
处理单聊消息
**/
func singleMsg(protocol *bean.UdpProtocol) {
	singleMsg := &bean.SingleMsg{}
	proto.Unmarshal(protocol.GetProtocolContent(), singleMsg)
	fmt.Println("传过来的消息", singleMsg)
	syncKey, err := redis.IncrUserSynckey(singleMsg.GetToUserId())
	if err != nil {
		fmt.Println("syncKey自增失败", err)
		msgSendRes(singleMsg.GetFromUserId(), protocol.GetFromaddress(), singleMsg.GetMsgId(), false)
		return
	}
	singleMsg.SyncKey = syncKey
	err = redis.SetUserSingleMsg(singleMsg)
	if err != nil {
		fmt.Println("保存消息失败", err)
		msgSendRes(singleMsg.GetFromUserId(), protocol.GetFromaddress(), singleMsg.GetMsgId(), false)
		return
	}
	//给发送者发响应
	msgSendRes(singleMsg.GetFromUserId(), protocol.GetFromaddress(), singleMsg.GetMsgId(), true)
	//给接收者发通知
	msgInform(singleMsg.GetToUserId(), syncKey)
}

/**
消息发送响应
**/
func msgSendRes(userId int64, addreess string, msgId string, flag bool) {

	udpclient, err := net.Dial("udp", addreess)
	defer udpclient.Close()
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	msgSendRes := &bean.MsgSendRes{Flag: true, MsgId: msgId}
	content, _ := proto.Marshal(msgSendRes)
	protocolRes := &bean.UdpProtocol{ToUuserId: userId, ProtocolContent: content, ProtocolType: int32(bean.ProtocolTypeEnum_MSG_SEND_RES)}
	b, _ := proto.Marshal(protocolRes)
	udpclient.Write(b)
	udpclient.Close()
}

/**
新消息通知
**/
func msgInform(userId int64, syncKey int64) {
	u := redis.GetOnlineUser(userId)
	if u == nil {
		fmt.Println("用户不在线", userId)
		return
	}
	udpclient, err := net.Dial("udp", u.Onlineaddr)
	defer udpclient.Close()
	if err != nil {
		fmt.Println("udp客户端建立失败", err)
	}
	toMsg := &bean.MsgInform{SyncKey: syncKey, UserId: userId}
	b, _ := proto.Marshal(toMsg)
	protocolRes := &bean.UdpProtocol{ToUuserId: userId, ProtocolContent: b, ProtocolType: int32(bean.ProtocolTypeEnum_MSG_INFORM)}
	t, _ := proto.Marshal(protocolRes)
	udpclient.Write(t)
	udpclient.Close()
}
