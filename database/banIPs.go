package database

import (
	"gorm.io/gorm"
	"time"
)

// 黑名单模型等数据库相关操作

// BanIPs 黑名单表的model
type BanIPs struct {
	Id         int    `gorm:"column:id"`
	Ip         string `gorm:"column:ip"`
	CreateTime int    `gorm:"column:create_time"`
	DeadTime   int    `gorm:"column:dead_time"`
}

// CreateBanIP 创建一个新的黑名单IP记录
func CreateBanIP(db *gorm.DB, ip string, DeadTime int) (*BanIPs, error) {
	currTime := int(time.Now().Unix())
	banIP := &BanIPs{Ip: ip, CreateTime: currTime, DeadTime: DeadTime}
	if err := db.Create(banIP).Error; err != nil {
		return nil, err
	}
	return banIP, nil
}

// GetBanIPById 根据Id获取黑名单IP记录
func GetBanIPById(db *gorm.DB, id int) (*BanIPs, error) {
	var banIP BanIPs
	if err := db.Where("id = ?", id).First(&banIP).Error; err != nil {
		return nil, err
	}
	return &banIP, nil
}

// GetBanIPByIP 根据IP获取黑名单记录
func GetBanIPByIP(db *gorm.DB, ip string) (*BanIPs, error) {
	var banIP BanIPs
	if err := db.Where("ip = ?", ip).First(&banIP).Error; err != nil {
		return nil, err
	} else {
		return &banIP, nil
	}
}

// DeleteBanIP 删除黑名单IP记录
func DeleteBanIP(db *gorm.DB, id int) error {
	if err := db.Where("id = ?", id).Delete(BanIPs{}).Error; err != nil {
		return err
	}
	return nil
}
