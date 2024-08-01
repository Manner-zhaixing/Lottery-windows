package main

import (
	"context"
	"fmt"
	"time"
)

const (
	WorkUseTime = 500 * time.Millisecond
	Timeout     = 100 * time.Millisecond
)

// 模拟一个耗时较长的任务
func LongTimeWork() int {
	time.Sleep(WorkUseTime)
	return 888
}

// 模拟一个接口处理函数
func Handle1() int {
	deadline := make(chan struct{}, 1)
	workDone := make(chan int, 1)
	go func() { //把要控制超时的函数放到一个协程里
		n := LongTimeWork()
		workDone <- n
	}()
	go func() { //把要控制超时的函数放到一个协程里
		time.Sleep(Timeout)
		// deadline <- struct{}{}
		close(deadline)
	}()
	select { //下面的case只执行最早到来的那一个
	case n := <-workDone:
		fmt.Println("LongTimeWork return")
		return n
	case <-deadline:
		fmt.Println("LongTimeWork timeout")
		return 0
	}
}

// 模拟一个接口处理函数
func Handle2() int {
	workDone := make(chan int, 1)
	go func() { //把要控制超时的函数放到一个协程里
		n := LongTimeWork()
		workDone <- n
	}()
	select { //下面的case只执行最早到来的那一个
	case n := <-workDone:
		fmt.Println("LongTimeWork return")
		return n
	case <-time.After(Timeout):
		fmt.Println("LongTimeWork timeout")
		return 0
	}
}

// 模拟一个接口处理函数
func Handle3() int {
	//通过显式sleep再调用cancle()来实现对函数的超时控制
	ctx, cancel := context.WithCancel(context.Background())

	workDone := make(chan int, 1)
	go func() { //把要控制超时的函数放到一个协程里
		n := LongTimeWork()
		workDone <- n
	}()

	go func() {
		//100毫秒后调用cancel()，关闭ctx.Done()
		time.Sleep(Timeout)
		cancel()
	}()

	select { //下面的case只执行最早到来的那一个
	case n := <-workDone:
		fmt.Println("LongTimeWork return")
		return n
	case <-ctx.Done(): //ctx.Done()是一个管道，调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
		return 0
	}
}

// 模拟一个接口处理函数
func Handle4() int {
	//借助于带超时的context来实现对函数的超时控制
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel() //纯粹出于良好习惯，函数退出前调用cancel()
	workDone := make(chan int, 1)
	go func() { //把要控制超时的函数放到一个协程里
		n := LongTimeWork()
		workDone <- n
	}()
	select { //下面的case只执行最早到来的那一个
	case n := <-workDone:
		fmt.Println("LongTimeWork return")
		return n
	case <-ctx.Done(): //ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		fmt.Println("LongTimeWork timeout")
		return 0
	}
}

func main17() {
	Handle1()
	Handle2()
	Handle3()
	Handle4()
}
