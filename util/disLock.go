package util

import (
	"context"
	"gift/database"
)

// 分布式锁(Redis实现)工具函数，根据用户id锁定

// GetDisLockKey 用于根据uid获取分布式锁的key值
func GetDisLockKey(uid string) string {
	key := "lucky_lock_" + uid
	return key
}

// DisLock 用于根据uid获取分布式锁
func DisLock(uid string) bool {
	ctx := context.Background()
	key := GetDisLockKey(uid)
	// NX参数：如果不存在，则设置分布式锁
	cmd, _ := database.GetRedisClient().Do(ctx, "SET", key, 1, "EX", 3, "NX").Result()
	if cmd == "OK" {
		return true
	} else {
		return false
	}
}

// DisUnLock 用于根据uid释放分布式锁
func DisUnLock(uid string) bool {
	ctx := context.Background()
	key := GetDisLockKey(uid)
	cmd, _ := database.GetRedisClient().Del(ctx, key).Result()
	if cmd == 1 {
		return true
	} else {
		return false
	}
}
