package main

import (
	"fmt"
	"runtime"
	"time"
)

type GoroutineLimiter struct {
	limit int //缓冲长度为limit，运行的协程不会超过这个值
	ch    chan struct{}
}

func NewGoroutineLimiter(n int) *GoroutineLimiter {
	return &GoroutineLimiter{
		limit: n,
		ch:    make(chan struct{}, n),
	}
}

func (g *GoroutineLimiter) Run(f func()) { //函数作这参数
	g.ch <- struct{}{} //创建子协程前往管道里send一个数据
	go func() {
		f()
		<-g.ch //子协程退出时从管理里取出一个数据
	}()
}

func main16() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	go func() {
		//每隔1秒打印一次协程数量
		for {
			<-ticker.C
			fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine())
		}
	}()

	limiter := NewGoroutineLimiter(100) //限制协程数为100
	work := func() {                    //函数变量
		time.Sleep(10 * time.Second)
	}
	for i := 0; i < 10000; i++ {
		// go work()
		limiter.Run(work) //不停地通过Run创建子协程
	}
	time.Sleep(10 * time.Second)
}
