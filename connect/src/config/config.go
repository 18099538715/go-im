package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

var udpIp, tcpIp, remoteUdpId, udpPort, tcpPort string
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
	udpIp = constantMap["udpIp"]
	udpPort = constantMap["udpPort"]
	tcpIp = constantMap["tcpIp"]
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
func GetUdpIp() string {
	return udpIp
}
func GetTcpPort() string {
	return tcpPort
}
func GetTcpIp() string {
	return tcpIp
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
