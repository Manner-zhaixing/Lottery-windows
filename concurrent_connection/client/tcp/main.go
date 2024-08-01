package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	conn net.Conn
)

func InitConn() {
	var err error
	//请求建立连接
	conn, err = net.DialTimeout("tcp", "localhost:5678", time.Duration(1*time.Second))
	if err != nil {
		panic(err)
	}
}

func DealSignal() {
	//按下Ctrl+C时，关闭连接
	var sigChan = make(chan os.Signal, 10)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-sigChan
	fmt.Printf("receive signal %d\n", sig)
	conn.Close()
	os.Exit(0)
}

func SendAndRceive() {
	//发送数据
	conn.Write([]byte("大乔乔"))
	// conn.Write([]byte("大乔乔"))
	// conn.Write([]byte("大乔乔"))
	//接收数据
	response := make([]byte, 256)
	n, _ := conn.Read(response)
	fmt.Println(string(response[:n]))
}

func main() {
	InitConn()
	go DealSignal()

	// SendAndRceive()

	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			SendAndRceive()
		}()
	}
	wg.Wait()
}

// go run ./concurrent_connection/client/tcp
