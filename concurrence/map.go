package main

import (
	"fmt"
	"sync"
)

// 原生map支持并发读，但不支持并发写或并发读写
var mp = make(map[int]int, 1000)
var smp = sync.Map{}

func readMap() {
	for i := 0; i < 1000; i++ {
		lock.RLock()
		_ = mp[10]
		lock.RUnlock()
		smp.Load(10)
	}
}

func writeMap() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		mp[5] = 5
		lock.Unlock()
		smp.Store(5, 5)
	}
}

// 用sync.Map来计数会存在问题
func mapInc(mp *sync.Map, key int) { //注意，必须传sync.Map的指针（要修改结构体，必须传指针）
	if oldValue, exists := mp.Load(key); exists {
		mp.Store(key, oldValue.(int)+1)
	} else {
		mp.Store(key, 1)
	}
}

func main14() {
	// // go readMap()
	// go readMap()
	// go writeMap()
	// // go writeMap()
	// time.Sleep(1 * time.Second)

	const P = 1000
	const key = 8
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			mapInc(&smp, key)
		}()
	}
	wg.Wait()

	value, _ := smp.Load(key)
	fmt.Println(value)

}
