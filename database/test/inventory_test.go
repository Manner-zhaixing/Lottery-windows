package test

import (
	"fmt"
	"gift/database"
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetAllGiftInventory(t *testing.T) {
	database.InitGiftInventory()
	gifts := database.GetAllGiftInventory()
	if len(gifts) == 0 {
		t.Fail()
	} else {
		for _, gift := range gifts {
			fmt.Printf("%d\t%d\n", gift.Id, gift.Count)
		}
	}
}

// 测试在并发减库存的情况下，逻辑不会出错
func TestReduceInventory(t *testing.T) {
	database.InitGiftInventory()
	const GiftId = 16 //此商品对应的库存为1000
	const P = 10      //开10个并发，减库存
	var success, fail int32
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 110; j++ {
				if err := database.ReduceInventory(GiftId); err == nil {
					atomic.AddInt32(&success, 1)
				} else {
					atomic.AddInt32(&fail, 1)
				}
			}
		}()
	}
	wg.Wait()
	if success != 1000 {
		t.Fail()
	}
	fmt.Printf("fail %d\n", fail)
}

// go test -v .\database\test\ -run=^TestGetAllGiftInventory$ -count=1
// go test -v .\database\test\ -run=^TestReduceInventory$ -count=1
