package database

import (
	"gift/util"
)

type Order struct {
	Id          int    `gorm:"column:id;primaryKey"`
	GiftId      int    `gorm:"column:gift_id"`
	UserId      int    `gorm:"column:user_id"`
	GiftName    string `gorm:"column:gift_name"`
	GiftPicture string `gorm:"column:gift_picture"`
	GiftPrice   int    `gorm:"column:gift_price"`
}

// CreateOrder 写入一条订单记录
func CreateOrder(order Order) int {
	db := GetGiftDBConnection()
	orderMsg := Order{
		GiftId:      order.GiftId,
		UserId:      order.UserId,
		GiftName:    order.GiftName,
		GiftPicture: order.GiftPicture,
		GiftPrice:   order.GiftPrice,
	}
	if err := db.Create(&orderMsg).Error; err != nil {
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
