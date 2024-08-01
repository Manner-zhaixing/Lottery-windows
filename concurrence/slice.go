package main

import (
	"fmt"
	"sync"
)

func main() {
	const LEN = 10          //切片的长度
	arr := make([]int, LEN) //初始化后，切片里包含LEN个0
	// const P = 2             //P个协程并行写arr
	const LOOP = 100 //每个协程把arr遍历LOOP次

	wg := sync.WaitGroup{}
	wg.Add(2)
	// for i := 0; i < P; i++ {
	go func() {
		defer wg.Done()
		for j := 0; j < LOOP; j++ {
			for index := 0; index < LEN; index++ {
				if index%10 == 1 {
					arr[index]++
				}
			}
		}
	}()
	// }
	// for i := 0; i < P; i++ {
	go func() {
		defer wg.Done()
		for j := 0; j < LOOP; j++ {
			for index := 0; index < LEN; index++ {
				if index > LEN/2 {
					arr[index]++
				}
			}
		}
	}()
	// }
	wg.Wait()

	sum := 0
	for _, ele := range arr {
		sum += ele
	}
	fmt.Println(sum)
}
