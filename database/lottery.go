package database

import (
	"fmt"
	"gift/config"
	"math/rand"
	"time"
)

// LotteryWeightedRandom 抽奖算法，权重需要手动传入,
func LotteryWeightedRandom(gifts []*Gift) (int, *Gift) {

	// 使用随机数生成器,基于当前时间生成随机数,生成一个介于0和总权重之间的随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWeight := r.Intn(config.TotalWeight)
	fmt.Println("随机数：", randomWeight)
	for _, v := range gifts {
		if v.MinWeight <= randomWeight && randomWeight < v.MaxWeight {
			return v.Id, v
		}
	}
	return 1, nil
}
