package main

import "gift/database"

func main4() {
	ch := make(chan database.Gift, 100)

	go func() {
		for i := 0; i < 10000; i++ {
			ch <- database.Gift{}
		}
		close(ch)
	}()

	for {
		gift, ok := <-ch
		if !ok { // ok==false 意味羞： 1.channel已空；2.channel已被关闭
			break
		}
		//把gift写入redis
		_ = gift
	}
}
