package main

import (
	"fmt"
	"math/rand"
	"time"
)

func listenMultiWay() {
	ch1 := make(chan int, 1000)
	ch2 := make(chan byte, 1000)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch1 <- rand.Int()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch2 <- byte(rand.Int())
		}
	}()

AB:
	for {
		select { //同时监听多个channel，谁先有数据先执行哪个case，如果同时有数据则随机选一个case执行
		case v1 := <-ch1:
			fmt.Printf("%d\n", v1)
		case v2 := <-ch2:
			fmt.Printf("%c\n", v2)
			if v2 < 40 {
				break AB
			}
		default:
			fmt.Println("default")
		}
	}

	select {
	case v1 := <-ch1:
		fmt.Printf("%d\n", v1)
	default:
	}
}

func main2() {
	listenMultiWay()
}

func readChanWithNonblock() {
	dataCh := make(chan int)

	go func() {
		for i := 0; i < 100; i++ {
			dataCh <- i
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < 100; i++ {
		select {
		case data := <-dataCh: //阻塞，等dataCh里有数据
			fmt.Println(data)
		default:
		}

		fmt.Println("YES!")
		time.Sleep(500 * time.Millisecond)
	}
}

type Task int //检索、添加、删除
type SearchEngine struct{}

func (SearchEngine) Do(Task) {}

func serve() {
	engine := new(SearchEngine) //初始化搜索引擎

	taskCh := make(chan Task, 100)
	stopCh := make(chan struct{}, 1)

	// 任务是异步执行的，先添加到channel里
	go func() {
		for {
			taskCh <- 1
		}
	}()

	// 发送stop信号
	go func() {
		time.Sleep(100 * time.Millisecond)
		stopCh <- struct{}{}
		engine = nil
	}()

	for {
		//select块会一直等待，直到有一个case解除阻塞
		select {
		case task := <-taskCh:
			if engine != nil {
				engine.Do(task) //异步执行任务
			}
		case <-stopCh:
			return //接收到stop信号时退出serve()函数
		}
	}
}
