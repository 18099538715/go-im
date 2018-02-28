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
func IncrUserSynckey(userId int64) (int64, error) {
	return redis.Int64(cluster.Do("INCR", "synckey_"+strconv.FormatInt(userId, 10)))
}

func GetCurrentUserSynckey(userId int64) (int64, error) {
	return redis.Int64(cluster.Do("GET", "currentSynckey_"+strconv.FormatInt(userId, 10)))
}
func SetCurrentUserSynckey(userId int64, syncKey int64) error {
	_, err := cluster.Do("SET", "currentSynckey_"+strconv.FormatInt(userId, 10), syncKey)
	return err
}
func GetUserMsgs(userId int64, syncKeys []int64) ([]*bean.SingleMsg, error) {

	msgs, err := redis.Values(cluster.Do("HMGET", "usermsgs_"+strconv.FormatInt(userId, 10), syncKeys))
	if err != nil {
		return nil, err

	}
	singleMsgs := make([]*bean.SingleMsg, len(msgs))
	for index, singleMsg := range msgs {
		s := &bean.SingleMsg{}
		err := proto.Unmarshal(singleMsg.([]byte), s)
		if err != nil {
			return nil, err
		}
		singleMsgs[index] = s

	}
	return singleMsgs, nil

}
func SetUserSingleMsg(msg *bean.SingleMsg) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = cluster.Do("HSET", "usermsgs_"+strconv.FormatInt(msg.ToUserId, 10), msg.SyncKey, b)
	return err
}
