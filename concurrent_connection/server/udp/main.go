package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	udpAddr, _ := net.ResolveUDPAddr("udp", "localhost:5678")
	conn, _ := net.ListenUDP("udp", udpAddr) //UDP不需要创建连接，所以不需要像TCP那样通过Accept()创建连接，这里的conn是个假连接。多个client可以共用这一个虚拟连接

	for {
		//读取连接上发来的请求数据
		request := make([]byte, 256)
		n, remoteAddr, err := conn.ReadFromUDP(request) //UDP是面向报文的，一次Read读出来的就是一个业务报文
		if err != nil {
			if err == io.EOF {
				break //如果对方关闭了连接，则
			} else {
				fmt.Println(err)
			}
		}
		//把响应返回给对方
		response := "hello " + string(request[:n])
		_, err = conn.WriteToUDP([]byte(response), remoteAddr)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
			}
		}
	}
}

// go run ./concurrent_connection/server/udp
