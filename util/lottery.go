package util

import (
	"gift/config"
	"gift/database"
	"math/rand"
	"time"
)

// Lottery 抽奖。给定每个奖品被抽中的概率（无需要做归一化，但概率必须大于0），返回被抽中的奖品下标
//func Lottery(probs []float64) int {
//	if len(probs) == 0 {
//		return -1
//	}
//	sum := 0.0
//	acc := make([]float64, 0, len(probs)) //累积概率
//	for _, prob := range probs {
//		sum += prob
//		acc = append(acc, sum)
//	}
//
//	// 获取(0,sum] 随机数，并看r落在acc中的哪一个位置
//	r := rand.Float64() * sum
//	index := BinarySearch(acc, r)
//	return index
//}
//
//// BinarySearch 二分法查找数组中>=target的最小的元素下标。arr是单调递增的(里面不能存在重复的元素)，如果target比arr的最后一个元素还大，则返回最后一个元素的下标
//func BinarySearch(arr []float64, target float64) int {
//	if len(arr) == 0 {
//		return -1
//	}
//	begin, end := 0, len(arr)-1
//
//	for {
//		//不论len(arr)为多少，都适用如下2个if
//		if target <= arr[begin] {
//			return begin
//		}
//		if target > arr[end] {
//			return end + 1
//		}
//
//		//len(arr)=2时，适用如下这个if
//		if begin == end-1 {
//			return end
//		}
//
//		//len(arr)>=3时，适用如下流程
//		middle := (begin + end) / 2
//		if arr[middle] > target {
//			end = middle
//		} else if arr[middle] < target {
//			begin = middle
//		} else {
//			return middle
//		}
//	}
//}
//
//// LotteryRandom 抽奖算法，随机抽奖
//func LotteryRandom(probs []float64) int {
//	giftLength := len(probs)
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	index := r.Intn(giftLength)
//	return index
//}

// LotteryWeightedRandom 抽奖算法，权重需要手动传入,
func LotteryWeightedRandom(gifts []*database.Gift) (int, *database.Gift) {

	// 使用随机数生成器,基于当前时间生成随机数,生成一个介于0和总权重之间的随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWeight := r.Intn(config.TotalWeight)
	for _, v := range gifts {
		if v.MinWeight <= randomWeight && randomWeight < v.MaxWeight {
			return v.Id, v
		}
	}
	return 1, nil
}
