package main

import (
	"fmt"
	"sync"
	"time"
)

var lock sync.RWMutex

// 读、写锁互斥
// 读锁 之间互相不排斥
// 写锁 之间互相排斥

func main11() {
	lock.RLock() //上读锁
	// lock.RUnlock() //释放读锁

	go func() {
		lock.Lock() //上写锁
		fmt.Println("上锁成功")
	}()

	time.Sleep(3 * time.Second)
}
