package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

var remoteUdpId, udpPort, tcpPort string
var connTimeOut, reomteUdpPort int64

func init() {
	var constantMap map[string]string = make(map[string]string)
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = json.Unmarshal(jsonFile, &constantMap)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
		return
	}
	udpPort = constantMap["udpPort"]
	tcpPort = constantMap["tcpPort"]
	remoteUdpId = constantMap["remoteUdpId"]
	connTimeOut, err = strconv.ParseInt(constantMap["connTimeOut"], 10, 64)
	if err != nil {
		fmt.Println("获取超时时间失败")
	}
	reomteUdpPort, err = strconv.ParseInt(constantMap["reomteUdpPort"], 10, 64)
	if err != nil {
		fmt.Println("获取远程udp端口失败")
	}
}
func GetUdpPort() string {
	return udpPort
}

func GetTcpPort() string {
	return tcpPort
}

func GetRemoteUdpPort() int {
	return int(reomteUdpPort)
}
func GetRemoteUdpIp() string {
	return remoteUdpId
}
func GetConnTimeOut() int64 {
	return connTimeOut
}
