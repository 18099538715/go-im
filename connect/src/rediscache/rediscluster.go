package rediscache

import (
	"bean"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chasex/redis-go-cluster"
)

var cluster *redis.Cluster

func init() {
	var err error
	cluster, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   []string{"118.89.182.47:5000", "118.89.182.47:5001", "118.89.182.47:5002"},
			ConnTimeout:  10 * time.Second,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
	if err != nil {
		fmt.Println("redis集群初始化错误", err)
	}
}
func SetOnlineUser(userInfo *bean.UserInfo) (interface{}, error) {
	a, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println("存储登录用户出错")
	}
	return cluster.Do("SET", "onlineinfo_"+strconv.FormatInt(userInfo.UserId, 10), a)
}
func DelOnlineUser(userId int64) (interface{}, error) {
	return cluster.Do("DEL", "onlineinfo_"+strconv.FormatInt(userId, 10))
}
func Test() {
	_, err := cluster.Do("GET", "onlineinfo_", "aaa")
	if err != nil {
		fmt.Println(err)
	}
}
