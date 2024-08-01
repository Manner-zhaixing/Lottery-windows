package database

import (
	"gift/util"

	"gorm.io/gorm"
)

const EMPTY_GIFT = 1 //空奖品（“谢谢参与”）的ID

type Gift struct {
	Id      int    `gorm:"column:id;primaryKey"`
	Name    string `gorm:"column:name"`
	Price   int    `gorm:"column:price"`
	Picture string `gorm:"column:picture"`
	Count   int    `gorm:"column:count"`
}

func (Gift) TableName() string {
	return "inventory"
}

var (
	_all_gift_field = util.GetGormFields(Gift{})
)

// GetAllGiftsV1 把gift表里的数据全部取出来。当数量不多时可以直接select * from table
func GetAllGiftsV1() []*Gift {
	db := GetGiftDBConnection()
	var gifts []*Gift
	err := db.Select(_all_gift_field).Find(&gifts).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			util.LogRus.Errorf("read table %s failed: %s", Gift{}.TableName(), err)
		}
	}
	return gifts
}

// GetAllGiftsV2 千万级以上大表遍历方案。
func GetAllGiftsV2(ch chan<- Gift) {
	db := GetGiftDBConnection()
	const PAGE_SIZE = 500
	maxid := 0
	for {
		var gifts []Gift
		err := db.Select(_all_gift_field).Where("id>?", maxid).Limit(PAGE_SIZE).Find(&gifts).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				util.LogRus.Errorf("read table %s failed: %s", Gift{}.TableName(), err)
			}
			break
		}
		if len(gifts) == 0 {
			break
		}
		for _, gift := range gifts {
			if gift.Id > maxid {
				maxid = gift.Id
			}
			ch <- gift
		}
	}
	close(ch) //close后就不允许再往channel里添加元素了
}
