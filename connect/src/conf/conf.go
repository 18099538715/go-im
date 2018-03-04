package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var ConfMap map[string]string = make(map[string]string)

func init() {

	jsonFile, err := ioutil.ReadFile("conf.json")
	if err != nil {
		fmt.Printf("yamlFile.Get err   #%v ", err)
	}
	err = json.Unmarshal(jsonFile, &ConfMap)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
}
