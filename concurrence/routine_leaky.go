package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 模拟：依赖一个外部的RPC接口
func rpc() int {
	time.Sleep(5 * time.Millisecond)
	return 888
}

// http接口
func homeHandler(ctx *gin.Context) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	workDone := make(chan int)
	go func() {
		n := rpc() // 去调用其他的接口，但是会对它进行超时控制
		workDone <- n
	}()
	select {
	case n := <-workDone:
		ctx.String(http.StatusOK, strconv.Itoa(n))
	case <-timeoutCtx.Done():
		ctx.String(http.StatusInternalServerError, strconv.Itoa(0))
	}
}

func main18() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		//每隔1秒打印一次协程数量
		for {
			<-ticker.C
			fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
		}
	}()

	//直接在游览器中访问http://127.0.0.1:3456/debug/pprof，里面有一项是goroutine。在package net/http/pprof的init()函数里指定了/debug/pprof的Handler
	go http.ListenAndServe("127.0.0.1:3456", nil) //在线prof

	gin.DefaultWriter = io.Discard
	engine := gin.Default()
	engine.GET("/", homeHandler)
	engine.Run("127.0.0.1:5678")
}
