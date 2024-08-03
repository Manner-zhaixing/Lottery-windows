package database

import (
	"context"
	"gift/config"
	"gift/util"
)

// 该文件提供拦截功能的redis等服务

// FlagBanIPsRedis 判断ip是否在黑名单中
func FlagBanIPsRedis(ip string) (bool, error) {
	ctx := context.Background()
	redisClient := GetRedisClient()
	flag, err := redisClient.SIsMember(ctx, config.BanIPsSetRedis, ip).Result()
	if err != nil {
		util.LogRus.Errorf("[database] FlagBanIPsRedis error: %v", err)
		return false, err
	}
	return flag, nil
}

// AddBanIPsRedis 将ip加入黑名单
func AddBanIPsRedis(ip string) error {
	redisClient := GetRedisClient()
	ctx := context.Background()
	err := redisClient.SAdd(ctx, config.BanIPsSetRedis, ip).Err()
	if err != nil {
		util.LogRus.Errorf("[database] AddBanIPsRedis error: %v", err)
		return err
	}
	// 设置黑名单的过期时间
	err = redisClient.Expire(ctx, ip, config.BanIPsExpireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// DelBanIPsRedis 将ip从黑名单中移除
func DelBanIPsRedis(ip string) error {
	redisClient := GetRedisClient()
	ctx := context.Background()
	err := redisClient.SRem(ctx, config.BanIPsSetRedis, ip).Err()
	if err != nil {
		util.LogRus.Errorf("[database] DelBanIPsRedis error: %v", err)
		return err
	}
	return nil
}

// FlagBanUsersRedis 判断user id是否在黑名单中
func FlagBanUsersRedis(userID string) (bool, error) {
	ctx := context.Background()
	redisClient := GetRedisClient()
	flag, err := redisClient.SIsMember(ctx, config.BanUsersSetRedis, userID).Result()
	if err != nil {
		util.LogRus.Errorf("[redis] Flag BanUsersRedis error: %v", err)
		return false, err
	}
	return flag, nil
}

// AddBanUsersRedis 将userID加入黑名单
func AddBanUsersRedis(userID string) error {
	redisClient := GetRedisClient()
	ctx := context.Background()
	err := redisClient.SAdd(ctx, config.BanUsersSetRedis, userID).Err()
	if err != nil {
		util.LogRus.Errorf("[redis] Add BanUsersRedis error: %v", err)
		return err
	}
	// 设置黑名单的过期时间
	err = redisClient.Expire(ctx, userID, config.BanUsersExpireTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// DelBanUsersRedis 将ip从黑名单中移除
func DelBanUsersRedis(userID string) error {
	redisClient := GetRedisClient()
	ctx := context.Background()
	err := redisClient.SRem(ctx, config.BanUsersSetRedis, userID).Err()
	if err != nil {
		util.LogRus.Errorf("[redis] Del BanUsersRedis error: %v", err)
		return err
	}
	return nil
}

//
