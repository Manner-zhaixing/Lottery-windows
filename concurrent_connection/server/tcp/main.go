package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	//解析地址
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", "localhost:5678")
	//开始监听
	listener, _ := net.ListenTCP("tcp4", tcpAddr)
	//无限循环。处理完一个client后，可以处理下一个
	for {
		//建立连接
		conn, _ := listener.Accept() //Accept()会阻塞，直到有人请求建立连接
		fmt.Printf("connect to %s\n", conn.RemoteAddr().String())
		for { //长连接，建立好连接后，进行多轮的收发
			//读取连接上发来的请求数据
			request := make([]byte, 256) //设定一个最大长度，防止flood attack
			n, err := conn.Read(request) //TCP是面向字节流的，不保证一次读出来的刚好就是一个业务报文
			if err != nil {
				if err == io.EOF {
					break //如果对方关闭了连接，则进行下一次Accept()
				} else {
					fmt.Println(err)
				}
			}
			//把响应返回给对方
			response := "hello " + string(request[:n])
			_, err = conn.Write([]byte(response))
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

// go run ./concurrent_connection/server/tcp
