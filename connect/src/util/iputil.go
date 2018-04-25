package util

import (
	"fmt"
	"net"
)

var localIp = ""

func GetLocalIP() string {
	if localIp != "" {
		return localIp
	}
	localIp = "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本地ip出错", err)
	}
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIp = ipnet.IP.String()
				return localIp
			}

		}
	}
	return localIp
}
