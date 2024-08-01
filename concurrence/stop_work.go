package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main1() {
	for {
		stopWork()
	}
}

func stopWork() {
	const Max = 100000
	const NumSenders = 100

	wgSenders := sync.WaitGroup{}
	wgSenders.Add(NumSenders)

	dataCh := make(chan int, 100) //sender和receiver传输数据使用的通道
	stopCh := make(chan struct{}, NumSenders)

	// senders
	for i := 0; i < NumSenders; i++ {
		go func() {
			defer wgSenders.Done()
			for {
				value := rand.Intn(Max)
				select {
				case <-stopCh: //接收到stop信号时，退出协程
					return
				case dataCh <- value: //生产数据，放入dataCh
				}
			}
		}()
	}

	// receiver
	go func() {
		for {
			value := <-dataCh
			if value == Max-1 { // 命中stop条件
				// for i := 0; i < NumSenders; i++ {
				// 	stopCh <- struct{}{} //发送stop信号，为了使写channel操作不阻塞，给stopCh设置一些容量
				// }
				close(stopCh)
				return //本协程立即能出
			}
		}
	}()

	wgSenders.Wait() //等所有sender结束
	fmt.Println("stop")
}
