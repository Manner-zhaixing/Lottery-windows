package util

import (
	"math/rand"
	"time"
)

// Lottery 抽奖。给定每个奖品被抽中的概率（无需要做归一化，但概率必须大于0），返回被抽中的奖品下标
func Lottery(probs []float64) int {
	if len(probs) == 0 {
		return -1
	}
	sum := 0.0
	acc := make([]float64, 0, len(probs)) //累积概率
	for _, prob := range probs {
		sum += prob
		acc = append(acc, sum)
	}

	// 获取(0,sum] 随机数，并看r落在acc中的哪一个位置
	r := rand.Float64() * sum
	index := BinarySearch(acc, r)
	return index
}

// BinarySearch 二分法查找数组中>=target的最小的元素下标。arr是单调递增的(里面不能存在重复的元素)，如果target比arr的最后一个元素还大，则返回最后一个元素的下标
func BinarySearch(arr []float64, target float64) int {
	if len(arr) == 0 {
		return -1
	}
	begin, end := 0, len(arr)-1

	for {
		//不论len(arr)为多少，都适用如下2个if
		if target <= arr[begin] {
			return begin
		}
		if target > arr[end] {
			return end + 1
		}

		//len(arr)=2时，适用如下这个if
		if begin == end-1 {
			return end
		}

		//len(arr)>=3时，适用如下流程
		middle := (begin + end) / 2
		if arr[middle] > target {
			end = middle
		} else if arr[middle] < target {
			begin = middle
		} else {
			return middle
		}
	}
}

// LotteryRandom 抽奖算法，随机抽奖
func LotteryRandom(probs []float64) int {
	giftLength := len(probs)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(giftLength)
	return index
}

// LotteryWeightedRandom 抽奖算法，根据奖品的权重随机选择一个奖品,权重需要手动传入,该算法适合结合其他用户信息，例如社交平台用户上线次数等
// @param probs 奖品库存
// @param weights 奖品权重
func LotteryWeightedRandom(probs []float64, weights []int) int {
	// 首先计算所有奖品的总权重
	var totalWeight int
	for _, w := range weights {
		totalWeight += w
	}

	if totalWeight == 0 {
		// 没有奖品
		return -1
	}

	// 使用随机数生成器,基于当前时间生成随机数,生成一个介于0和总权重之间的随机数
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWeight := r.Intn(totalWeight)

	// 累计权重，找到随机数落在哪个奖品的区间内
	var accWeight int
	for i, w := range weights {
		accWeight += w
		if randomWeight < accWeight {
			return i
		}
	}

	// 如果由于某种原因没有找到匹配的奖品，则返回最后一个奖品
	return -1
}
