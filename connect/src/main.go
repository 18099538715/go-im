package main

import (
	"server"
)

func main() {
	server.StartTcpServer()
	server.StartUdpClient()
	server.StartUdpServer()
	//这里是为了阻塞主程序一直执行
	c := make(chan int)
	<-c
}
