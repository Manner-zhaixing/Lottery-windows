package config

import (
	"time"
)

// 负责存储一些默认配置

var (
	// BanIPsSetRedis 负责存储黑名单ip的set名称，redis
	BanIPsSetRedis = "banListUserSet"
	// BanIPsExpireTime 黑名单ip过期时间
	BanIPsExpireTime = time.Hour
	// BanUsersSetRedis 负责存储黑名单用户的set名称，redis
	BanUsersSetRedis = "banListUserSet"
	// BanUsersExpireTime 黑名单用户过期时间
	BanUsersExpireTime = time.Hour
	// IPSumsDayLimitMax ip每日访问最大上限
	IPSumsDayLimitMax int64 = 1000000
	// UserSumsDayLimitMax 用户id每日访问最大上限
	UserSumsDayLimitMax int64 = 1000000
	// LotteryAlgorithm 抽奖算法种类
	LotteryAlgorithm = 1
	// GiftCountPrefix Mysql库存初始化存到Redis的key名字--只保存库存
	GiftCountPrefix = "gift_count_"
	// GiftGTypeMinMaxPrefix Mysql库存初始化存到Redis的key名字--只保存库存
	GiftGTypeMinMaxPrefix = "gift_gtype_min_max_"
	// TotalWeight 抽奖奖品总权重
	TotalWeight = 10000
	// DeadTime 黑名单的截止时间
	DeadTime = int(time.Now().Unix()) + 12*60*60
	// GuaranteePrefix 保底次数的前缀
	GuaranteePrefix = "guarantee:"
	// GuaranteeSum 保底次数，如果抽奖次数==保底次数，则抽取保底奖品
	GuaranteeSum = 100
	// GenerateRandomMaxUserID 生成随机userid的最大值
	GenerateRandomMaxUserID = 100000
)
