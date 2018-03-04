package rediscache

import (
	"bean"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chasex/redis-go-cluster"
	"github.com/golang/protobuf/proto"
)

var cluster *redis.Cluster

const msgExpireTime int64 = 60 * 60 * 24 * 7 //缓存消息7天过期
func init() {
	var err error
	cluster, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   []string{"118.89.182.47:5000", "118.89.182.47:5001", "118.89.182.47:5002", "118.89.182.47:5003", "118.89.182.47:5004", "118.89.182.47:5005"},
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
	if err != nil {
		fmt.Println("redis集群初始化错误", err)
	}
}
func GetOnlineUser(userId int64) *bean.UserInfo {
	reply, err := redis.Bytes(cluster.Do("GET", "onlineinfo_"+strconv.FormatInt(userId, 10)))
	if err != nil {
		fmt.Println("redis读取在线用户信息错误", err)
		return nil
	}
	u := &bean.UserInfo{}
	err = json.Unmarshal(reply, u)
	if err != nil {
		fmt.Println("在线用户信息转换失败", err)
		return nil
	}
	return u
}
func IncrUserSrlNo(userId int64) (int64, error) {
	return redis.Int64(cluster.Do("INCR", "srlNo_"+strconv.FormatInt(userId, 10)))
}

func GetCurrentUserSrlNo(userId int64) (int64, error) {
	return redis.Int64(cluster.Do("GET", "currentSrlNo_"+strconv.FormatInt(userId, 10)))
}
func SetCurrentUserSrlNo(userId int64, srlNo int64) error {
	_, err := cluster.Do("SET", "currentSrlNo_"+strconv.FormatInt(userId, 10), srlNo)
	return err
}
func GetUserMsgs(userId int64, srlNos []int64) ([]*bean.SingleMsg, error) {
	ids := make([]interface{}, len(srlNos))
	for index, _ := range ids {
		ids[index] = "{usermsgs_" + strconv.FormatInt(userId, 10) + "}_" + strconv.FormatInt(srlNos[index], 10)
	}
	fmt.Println(len(ids))
	msgs, err := redis.Values(cluster.Do("MGET", ids...))
	fmt.Println("用户消息", msgs)
	if err != nil {
		return nil, err

	}
	singleMsgs := make([]*bean.SingleMsg, len(msgs))
	for index, singleMsg := range msgs {
		s := &bean.SingleMsg{}
		if singleMsg != nil {
			err := proto.Unmarshal(singleMsg.([]byte), s)
			if err != nil {
				return nil, err
			}
			singleMsgs[index] = s
		}
	}
	return singleMsgs, nil

}
func SetUserSingleMsg(msg *bean.SingleMsg) error {

	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = cluster.Do("SETEX", "{usermsgs_"+strconv.FormatInt(msg.ToUserId, 10)+"}_"+strconv.FormatInt(msg.SrlNo, 10), msgExpireTime, b)
	fmt.Println("保存消息", err, msg)
	return err
}
