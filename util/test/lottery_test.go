package test

import (
	"fmt"
	"gift/util"
	"math/rand"
	"slices"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	const L = 100
	for c := 0; c < 30; c++ { //进行多轮严格的测试
		arr := make([]float64, 0, L)
		for i := 0; i < L; i++ {
			arr = append(arr, rand.Float64())
		}
		slices.Sort(arr) //排序
		var target float64

		//先测试2个越界的情况
		target = arr[0] - 1.0
		if util.BinarySearch(arr, target) != 0 {
			t.Fail()
		}
		target = arr[len(arr)-1] + 1.0
		if util.BinarySearch(arr, target) != len(arr) {
			t.Fail()
		}

		// 每个分割点，及2个分割点中间的值都测一下
		target = arr[0]
		if util.BinarySearch(arr, target) != 0 {
			t.Fail()
		}

		for i := 0; i < len(arr)-1; i++ {
			target = (arr[i] + arr[i+1]) / 2
			if util.BinarySearch(arr, target) != i+1 {
				t.Fail()
			}
			target = arr[i+1]
			if util.BinarySearch(arr, target) != i+1 {
				t.Fail()
			}
		}
	}
}

// 测试按比例抽奖
func TestLottery(t *testing.T) {
	probs := []float64{5, 2, 4} //指定各元素被抽中的概率
	countMap := make(map[int]float64, len(probs))
	for i := 0; i < len(probs); i++ {
		countMap[i] = 0
	}
	for i := 0; i < 10000; i++ {
		index := util.Lottery(probs)
		countMap[index] += 1
	}

	// 以下值应该很接近
	fmt.Println(countMap[0] / probs[0])
	fmt.Println(countMap[1] / probs[1])
	fmt.Println(countMap[2] / probs[2])
}

// go test -v .\util\test\ -run=^TestBinarySearch$ -count=1
// go test -v .\util\test\ -run=^TestLottery$ -count=1
