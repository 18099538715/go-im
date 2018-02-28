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
func SetOnlineUser(userInfo *bean.UserInfo) (interface{}, error) {
	a, _ := json.Marshal(userInfo)
	return cluster.Do("SET", "onlineinfo_"+strconv.FormatInt(userInfo.UserId, 10), a)
}
