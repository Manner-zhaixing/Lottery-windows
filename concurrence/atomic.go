package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var n int32 = 0

func inc() {
	// n++
	// atomic.AddInt32(&n, 1)
	lock.Lock()
	n++
	lock.Unlock()
}

func main10() {
	const P = 1000
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc()
		}()
	}
	wg.Wait()
	fmt.Printf("n=%d\n", atomic.LoadInt32(&n))
}
