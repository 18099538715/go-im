package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var udpIp, udpPort string

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
}
func GetUdpPort() string {
	return udpPort
}
func GetUdpIp() string {
	return udpIp
}
