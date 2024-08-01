package main

import (
	"fmt"
	"sync"
	"time"
)

var ch = make(chan int, 3)
var twg = sync.WaitGroup{}

func add2Ch() {
	defer twg.Done()
	for i := 0; i < 10; i++ {
		ch <- i
	}
	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	// close(ch)
}

func traveseChannel() {
	defer twg.Done()
	for ele := range ch { //遍历并取走管道中的元素
		fmt.Println(ele)
	}

	// for {
	// 	ele, ok := <-ch
	// 	if ok {
	// 		fmt.Println(ele)
	// 	} else {
	// 		break
	// 	}
	// }

	fmt.Println("bye bye")
}

func main9() {
	twg.Add(2)
	go add2Ch()
	go traveseChannel()
	go func() {
		// time.Sleep(time.Hour)
		<-ch
	}()

	// time.Sleep(3 * time.Second)
	twg.Wait()
}
