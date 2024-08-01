package test

import (
	"gift/database"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	userId, giftId := 3, 6
	orderId := database.CreateOrder(userId, giftId)
	if orderId <= 0 {
		t.Fail()
	}
}

// go test -v .\database\test\ -run=^TestCreateOrder$ -count=1
