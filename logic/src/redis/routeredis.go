package redis

import (
	"fmt"
	"time"

	"github.com/chasex/redis-go-cluster"
)

var Cluster *redis.Cluster

func init() {
	var err error
	Cluster, err = redis.NewCluster(
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
