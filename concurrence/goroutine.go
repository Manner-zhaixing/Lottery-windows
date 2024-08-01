package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	wg = sync.WaitGroup{} //计数为0
)

func init() {
	wg.Add(2) //计数为2
}

func parent() {
	defer wg.Done() //计数减1
	go child()
	for i := 'a'; i <= 'z'; i++ {
		fmt.Printf("%d\n", i)
		// time.Sleep(500 * time.Millisecond)
	}
}

func child() {
	defer wg.Done() //计数减1
	for i := 'a'; i <= 'z'; i++ {
		fmt.Printf("%c\n", i)
		time.Sleep(50 * time.Millisecond)
	}
}

func main5() { // runtime

	// go parent() //开启了一个子协程

	go func() { //匿名函数
		defer wg.Done()
		go child()
		for i := 'a'; i <= 'z'; i++ {
			fmt.Printf("%d\n", i)
			// time.Sleep(50 * time.Millisecond)
		}
	}()

	// go child()
	fmt.Println("main") //main 协程

	// time.Sleep(100 * time.Second)
	wg.Wait() //阻塞，直到计数减为0

	cpuN := runtime.NumCPU()
	fmt.Println("逻辑核数", cpuN)
	runtime.GOMAXPROCS(cpuN / 2) //限制go进程最多使用的核数

	const P = 1000000
	for i := 0; i < P; i++ {
		go time.Sleep(10 * time.Second)
	}
	fmt.Println("进程中存活的协程数", runtime.NumGoroutine())
}
