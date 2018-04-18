package rediscache

import (
	"bean"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/chasex/redis-go-cluster"
)

var cluster *redis.Cluster

func init() {
	var servers map[string]string = make(map[string]string)
	jsonFile, err := ioutil.ReadFile("redis.json")
	if err != nil {
		fmt.Printf("读取redis配置文件出错", err)
	}
	err = json.Unmarshal(jsonFile, &servers)
	if err != nil {
		fmt.Println("json redis配置信息出错", err)
	}
	serverList := strings.Split(servers["serverList"], ",")
	cluster, err = redis.NewCluster(
		&redis.Options{
			StartNodes:   serverList,
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
	fmt.Println(userInfo)
	a, err := json.Marshal(userInfo)
	if err != nil {
		fmt.Println("存储登录用户出错")
	}
	c, d := cluster.Do("SET", "onlineinfo_"+strconv.FormatInt(userInfo.UserId, 10), a)
	fmt.Println(c, d)
	return c, d
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
