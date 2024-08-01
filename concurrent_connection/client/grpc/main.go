package main

import (
	"context"
	"fmt"
	"gift/concurrent_connection/server/idl"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn *grpc.ClientConn
)

func InitConn() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	var err error
	conn, err = grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}

}

func Call() {
	//创建client
	client := idl.NewHelloClient(conn)
	request := idl.HelloRequest{Data: "大乔乔"}
	response, err := client.Hello(context.Background(), &request)
	if err == nil {
		fmt.Println(response.Data)
	} else {
		fmt.Println(err)
	}
}

func main() {
	InitConn()
	defer conn.Close()

	Call()

	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			Call()
		}()
	}
	wg.Wait()
}

// go run ./concurrent_connection/client/grpc
