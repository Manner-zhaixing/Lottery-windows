package database

import (
	"context"
	"errors"
	"gift/util"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// 该文件提供拦截功能的redis等服务,利用hash实现
// 思路是：简历hash类型，key为ip或者userid，value为1，为每个key过期时间，
// 每次请求时查询key值，如果不存在说明是这段时间内第一次访问，则可以抽奖，并将ip或者userid作为key，1作为value存入hash
// 如果查询到的value值存在，则判断value是否超过限定的一定时间内的请求次数，如果超过则是为恶意请求，则返回错误，否则继续抽奖，value+1

const userIDFrameHash = 2

// InitHashRedisTasks 定时任务，0点删除用户每日抽奖次数和ip每日抽奖次数，需要go协程开启
func InitHashRedisTasks() {
	c := cron.New()
	// 添加定时任务，每天午夜0点执行
	c.AddFunc("0 0 * * *", func() {
		resetUserSumsDayHashRedis()
		resetIPSumsDayHashRedis()
	})
	c.Start()
}

// resetUserSumsDayHashRedis 用于每天0点重置用户每日抽奖次数，redis中的hash结构
func resetUserSumsDayHashRedis() {
	ctx := context.Background()
	for i := 0; i < userIDFrameHash; i++ {
		key := "userid_Sums_Day_" + strconv.Itoa(i)
		GetRedisClient().Del(ctx, key)
	}
}

// IncrementUserSumsDay 原子性增加抽奖次数
func IncrementUserSumsDay(userID string) int64 {
	// 对userid进行散列，某几个ip对应在一个hash中
	userIDInt, _ := strconv.Atoi(userID)
	i := userIDInt % ipFrameHash
	key := "userid_Sums_Day_" + strconv.Itoa(i)
	ctx := context.Background()
	// hash增加次数
	cmd, err := GetRedisClient().HIncrBy(ctx, key, userID, 1).Result()
	if err != nil {
		util.LogRus.Errorf("userid add failed")
		return -1
	} else {
		return cmd
	}
}

// CheckBanUsers 检查用户是否在黑名单
func CheckBanUsers(userid int) bool {
	mysqlClient := GetGiftDBConnection()
	banUser, err := GetBanUserByUserID(mysqlClient, userid)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 不在黑名单
		return false
	}
	if banUser != nil {
		if banUser.LimitTime > int(time.Now().Unix()) {
			// 还在黑名单中，且未过期
			return true
		} else {
			// 在黑名单中，过期了
			return false
		}
	}
	return false
}

// 用于处理ip每日次数的工具函数
const ipFrameHash = 2

// resetIPSumsDayHashRedis 重置ip每日次数
func resetIPSumsDayHashRedis() {
	ctx := context.Background()
	for i := 0; i < ipFrameHash; i++ {
		key := "ip_Sums_Day_" + strconv.Itoa(i)
		GetRedisClient().Del(ctx, key)
	}
}

// IncrementIpSumsDay 增加ip每日次数
func IncrementIpSumsDay(ip string) int64 {
	// 对ip进行散列，某几个ip对应在一个hash中
	ipInt, _ := strconv.Atoi(ip)
	i := ipInt % ipFrameHash
	key := "ip_Sums_Day_" + strconv.Itoa(i)
	ctx := context.Background()
	// hash增加次数
	cmd, err := GetRedisClient().HIncrBy(ctx, key, ip, 1).Result()
	if err != nil {
		// 次数增加失败
		util.LogRus.Errorf("ipSumsDay IncrementIpSumsDay error: %v", err)
		return -1
	} else {
		return cmd
	}
}

// CheckBanIPs 检查ip是否在黑名单
func CheckBanIPs(ip string) bool {
	mysqlClient := GetGiftDBConnection()
	banIP, err := GetBanIPByIP(mysqlClient, ip)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// 不在黑名单
		return false
	}
	if banIP != nil {
		if banIP.LimitTime > int(time.Now().Unix()) {
			// 还在黑名单中，且未过期
			return true
		} else {
			// 在黑名单中，过期了
			return false
		}
	}
	return false
}
