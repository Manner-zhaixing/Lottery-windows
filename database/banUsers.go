package database

import (
	"gorm.io/gorm"
	"time"
)

// BanUsers 黑名单表的model
type BanUsers struct {
	Id         int `gorm:"column:id"`
	UserID     int `gorm:"column:userid"`
	CreateTime int `gorm:"column:create_time"`
	DeadTime   int `gorm:"column:dead_time"`
}

// CreateBanUser 创建一个新的用户黑名单记录
func CreateBanUser(db *gorm.DB, userID int, deadTime int) (*BanUsers, error) {
	currTime := int(time.Now().Unix())
	banUser := &BanUsers{UserID: userID, CreateTime: currTime, DeadTime: deadTime}
	if err := db.Create(banUser).Error; err != nil {
		return nil, err
	}
	return banUser, nil
}

// GetBanUserByID 根据ID获取用户黑名单记录
func GetBanUserByID(db *gorm.DB, id int) (*BanUsers, error) {
	var banUser BanUsers
	if err := db.Where("id = ?", id).First(&banUser).Error; err != nil {
		return nil, err
	}
	return &banUser, nil
}

// GetBanUserByUserID 根据UserID获取用户黑名单记录
func GetBanUserByUserID(db *gorm.DB, userID int) (*BanUsers, error) {
	var banUser BanUsers
	if err := db.Where("userid = ?", userID).First(&banUser).Error; err != nil {
		return nil, err
	}
	return &banUser, nil
}

// GetAllBanUsers 获取所有用户黑名单记录
func GetAllBanUsers(db *gorm.DB) ([]BanUsers, error) {
	var banUsers []BanUsers
	if err := db.Find(&banUsers).Error; err != nil {
		return nil, err
	}
	return banUsers, nil
}

// DeleteBanUser 删除用户黑名单记录
func DeleteBanUser(db *gorm.DB, id int) error {
	if err := db.Where("id = ?", id).Delete(&BanUsers{}).Error; err != nil {
		return err
	}
	return nil
}
