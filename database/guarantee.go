package database

import (
	"context"
	"gift/config"
	"gift/util"
	"strconv"
)

// 该文件存放关于抽奖保底的函数

// GuaranteeSumReset 保底次数重置
func GuaranteeSumReset(userID string) {
	// 拼接Redis中保底次数的key
	key := config.GuaranteePrefix + userID
	ctx := context.Background()
	err := GetRedisClient().Set(ctx, key, 0, 0).Err()
	if err != nil {
		util.LogRus.Errorf("用户保底抽奖次数重置失败")
	}
}

// GuaranteeSumIncrement 保底次数增加,redis的原子操作
func GuaranteeSumIncrement(userID string) int64 {
	// 拼接Redis中保底次数的key
	key := config.GuaranteePrefix + userID
	ctx := context.Background()
	n, err := GetRedisClient().Incr(ctx, key).Result()
	if err != nil {
		util.LogRus.Errorf("用户保底抽奖次数增加失败")
	}
	return n
}

// GuaranteeSumNow 获取目前的抽奖保底次数
func GuaranteeSumNow(userID string) int {
	ctx := context.Background()
	key := config.GuaranteePrefix + userID
	cmd, err := GetRedisClient().Get(ctx, key).Result()
	n, _ := strconv.Atoi(cmd)
	if err != nil {
		util.LogRus.Errorf("获取用户保底抽奖次数失败")
		return -1
	} else {
		return n
	}
}
