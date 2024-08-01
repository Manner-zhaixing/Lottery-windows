package test

import (
	"fmt"
	"gift/database"
	"gift/util"
	"testing"
)

func init() {
	util.InitLog("log")
}

func TestGetAllGiftsV1(t *testing.T) {
	gifts := database.GetAllGiftsV1()
	if len(gifts) == 0 {
		t.Fail()
	} else {
		for _, gift := range gifts {
			fmt.Printf("%+v\n", *gift)
		}
	}
}

func TestGetAllGiftsV2(t *testing.T) {
	giftCh := make(chan database.Gift, 100)
	go database.GetAllGiftsV2(giftCh)
	for {
		gift, ok := <-giftCh
		if !ok { //channel已经消费完了
			break
		}
		fmt.Printf("%+v\n", gift)
	}
}

// go test -v .\database\test\ -run=^TestGetAllGiftsV1$ -count=1
// go test -v .\database\test\ -run=^TestGetAllGiftsV2$ -count=1
