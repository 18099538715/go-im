package rediscache

import (
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
	fmt.Println(serverList)
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
	_, err = cluster.Do("SET", "a1aa", "cc")
	if err != nil {
		fmt.Println(err, "第一次set错误")
	}
	_, err = cluster.Do("SET", "a1aa", "cc")
	if err != nil {
		fmt.Println(err, "第一次set错误")
	}
	_, err = cluster.Do("SET", "aaa", "cc")
	if err != nil {
		fmt.Println(err, "第二次set错误")
	}
	_, err = cluster.Do("SET", "bbb", "cc")
	if err != nil {
		fmt.Println(err, "第三次set错误")
	}
	_, err = cluster.Do("SET", "ccc", "cc")
	if err != nil {
		fmt.Println(err, "第四次set错误")
	}
	fmt.Println("第二次成功")
}
func SetUserToken(userId int64, token string) (interface{}, error) {
	fmt.Println("往redis里写数据")
	return cluster.Do("SET", "usertoken_"+strconv.FormatInt(userId, 10), token)
}
