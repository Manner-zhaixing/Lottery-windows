package database

import (
	"gift/util"
)

type Order struct {
	Id     int
	GiftId int
	UserId int
}

// CreateOrder 写入一条订单记录
func CreateOrder(userid, giftid int) int {
	db := GetGiftDBConnection()
	order := Order{GiftId: giftid, UserId: userid}
	if err := db.Create(&order).Error; err != nil {
		util.LogRus.Errorf("create order failed: %s", err)
		return 0
	} else {
		util.LogRus.Debugf("create order id %d", order.Id)
		return order.Id
	}
}

// ClearOrders 清除全部订单记录
func ClearOrders() error {
	db := GetGiftDBConnection()
	return db.Where("id>0").Delete(Order{}).Error
}
