package test

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const url = "http://localhost:5678/lucky"
const P = 100 //模拟100个用户，在疯狂地点击“抽奖”

func TestLottery1(t *testing.T) {
	hitMap := sync.Map{}

	wg := sync.WaitGroup{}
	wg.Add(P)
	begin := time.Now()
	var totalCall int64    //记录接口总调用次数
	var totalUseTime int64 //接口调用耗时总和
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for {
				t1 := time.Now()
				resp, err := http.Get(url)
				atomic.AddInt64(&totalUseTime, time.Since(t1).Milliseconds())
				atomic.AddInt64(&totalCall, 1) //调用次数加1
				if err != nil {
					fmt.Println(err)
					break
				}
				bs, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					break
				}
				resp.Body.Close()
				giftId := string(bs)
				if giftId == "0" { //如果返回的奖品ID为0，说明已抽完
					break
				}
				if cnt, exists := hitMap.Load(giftId); exists {
					hitMap.Store(giftId, cnt.(int)+1) //多个协程同时执行？计数被覆盖，会比实际的小
				} else {
					hitMap.Store(giftId, 1) //多个协程同时执行？
				}
			}
		}()
	}
	wg.Wait()
	totalTime := int64(time.Since(begin).Seconds())
	if totalTime > 0 && totalCall > 0 {
		qps := totalCall / totalTime
		avgTime := totalUseTime / totalCall
		fmt.Printf("QPS %d, avg time %dms\n", qps, avgTime) //QPS 2800, avg time 38ms。这里的QPS比server实际上吞吐能力低，因为sync.Map效率比较低，读写hitMap消耗了比较多的时间，换成channel就会快很多

		total := 0
		hitMap.Range(func(giftId, count any) bool {
			fmt.Printf("%s\t%d\n", giftId, count.(int))
			total += count.(int)
			return true
		})
		fmt.Printf("共计%d件商品\n", total)
	}
}

func TestLottery2(t *testing.T) {
	hitMap := make(map[string]int, 10)
	giftCh := make(chan string, 10000)
	counterCh := make(chan struct{})

	go func() {
		for {
			giftId, ok := <-giftCh
			if !ok {
				break
			}
			if cnt, exists := hitMap[giftId]; exists {
				hitMap[giftId] = cnt + 1
			} else {
				hitMap[giftId] = 1
			}
		}
		counterCh <- struct{}{}
	}()

	wg := sync.WaitGroup{}
	wg.Add(P)
	begin := time.Now()
	var totalCall int64    //记录接口总调用次数
	var totalUseTime int64 //接口调用耗时总和
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for {
				t1 := time.Now()
				resp, err := http.Get(url)
				atomic.AddInt64(&totalUseTime, time.Since(t1).Milliseconds())
				atomic.AddInt64(&totalCall, 1) //调用次数加1
				if err != nil {
					fmt.Println(err)
					break
				}
				bs, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					break
				}
				resp.Body.Close()
				giftId := string(bs)
				if giftId == "0" { //如果返回的奖品ID为0，说明已抽完
					break
				}
				giftCh <- giftId //抽中一个gift，就放入channel
			}
		}()
	}
	wg.Wait()
	close(giftCh)
	<-counterCh //等hitMap准备好

	totalTime := int64(time.Since(begin).Seconds())
	if totalTime > 0 && totalCall > 0 {
		qps := totalCall / totalTime
		avgTime := totalUseTime / totalCall
		fmt.Printf("QPS %d, avg time %dms\n", qps, avgTime) //QPS 5600, avg time 31ms

		total := 0
		for giftId, count := range hitMap {
			fmt.Printf("%s\t%d\n", giftId, count)
			total += count
		}
		fmt.Printf("共计%d件商品\n", total)
	}
}

// go test -v .\handler\test\ -run=^TestLottery1$ -count=1
// go test -v .\handler\test\ -run=^TestLottery2$ -count=1
