package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error: ", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			break
		}

		// start a new goroutine to handle the new connection
		go HandleConn(conn)
	}

}
func HandleConn(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr()
	fmt.Println(remoteAddr.String())
	for {
		time.Sleep(2 * time.Second)
		conn.Write([]byte("connect success"))
	}

}
