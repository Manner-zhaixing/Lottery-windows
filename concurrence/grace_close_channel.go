package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

func main3() {
	for {
		closeCh()
	}
}

func closeCh() {
	const Max = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)
	wgSenders := sync.WaitGroup{}
	wgSenders.Add(NumSenders)

	dataCh := make(chan int, 100)  //sender和receiver传输数据使用的通道
	stopCh := make(chan struct{})  //由主持人向stopCh里发送数据，sender和receiver都监听此通道，通道关闭时sender和receiver立即终止
	toStop := make(chan string, 1) //sender和receiver都可以向toStop里发送数据，主持人监听此通道，通道里有数据时就关闭stopCh

	var stoppedBy string //记录stop请求最初是由谁发起的（sender或receiver）

	// 主持人
	go func() {
		stoppedBy = <-toStop //主持人监听toStop
		close(stopCh)        //toStop里有数据时就关闭stopCh
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			defer wgSenders.Done()
			for {
				value := rand.Intn(Max)
				if value == 0 { // 命中stop条件
					// 通知主持人，关闭stopCh
					select {
					case toStop <- "sender#" + id: //toStop容量为1，避免写toStop阻塞
					default:
					}
					// stoppedBy = "sender#" + id
					// close(stopCh)
					return //本协程立即能出
				}

				// 正常的业务逻辑，向dataCh里生产数据
				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()
			for {
				// 正常的业务逻辑，从dataCh里取出数据并消费
				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == Max-1 { // 命中stop条件
						// 通知主持人，关闭stopCh
						select {
						case toStop <- "receiver#" + id: //toStop容量为1，避免写toStop阻塞
						default:
						}
						// stoppedBy = "receiver#" + id
						// close(stopCh)
						return //本协程立即能出
					}
				}
			}
		}(strconv.Itoa(i))
	}

	wgReceivers.Wait() //等所有receiver结束
	wgSenders.Wait()   //等所有sender结束
	fmt.Println("stopped by", stoppedBy)
}
