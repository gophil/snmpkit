package main

import (
	"github.com/gophil/snmpkit/config"
)

func main() {

	//读取交换机的设备信息
	items, err := config.LoadSwitchFromFile("/Users/lihaoquan/Desktop/switch..")
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		println(item.Host)
	}
}
