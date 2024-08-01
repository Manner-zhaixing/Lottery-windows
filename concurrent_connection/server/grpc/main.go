package main

import (
	"context"
	"errors"
	"gift/concurrent_connection/server/idl"
	"net"

	"google.golang.org/grpc"
)

type MyHelloService struct {
	idl.UnimplementedHelloServer //必须继承UnimplementedHelloServer才实现了接口
}

func (s MyHelloService) Hello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	if request == nil {
		return nil, errors.New("empty request")
	}
	return &idl.HelloResponse{Data: "hello " + request.Data}, nil
}

func main() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	//创建服务
	server := grpc.NewServer()
	// 注册服务的具体实现
	idl.RegisterHelloServer(server, &MyHelloService{})
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run ./concurrent_connection/server/grpc
