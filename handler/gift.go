package handler

import (
	"gift/database"
	"gift/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllGifts 获取所有奖品信息，用于初始化轮盘
func GetAllGifts(ctx *gin.Context) {
	gifts := database.GetAllGiftsV1()
	if len(gifts) == 0 {
		ctx.JSON(http.StatusInternalServerError, nil)
	} else {
		//抹掉敏感信息
		for _, gift := range gifts {
			gift.Count = 0
		}
		ctx.JSON(http.StatusOK, gifts)
	}
}

// Lottery 抽奖
func Lottery(ctx *gin.Context) {
	// 使用的抽奖算法类别 [1-库存作为权重,2-随机,3-手动赋值权重,可以结合其他业务]
	var lotteryAlgorithm int
	lotteryAlgorithm = 1
	var index int
	for try := 0; try < 3; try++ { //最多重试10次
		// gifts := make([]*Gift, 0, len(keys))
		// Gift{Id: id, Count: count}
		gifts := database.GetAllGiftInventory() //获取所有奖品剩余的库存量Redis
		ids := make([]int, 0, len(gifts))
		probs := make([]float64, 0, len(gifts))
		for _, gift := range gifts {
			if gift.Count > 0 { //先确保redis返回的库存量大小0，因为抽奖算法Lottery不支持抽中概率为0的奖品
				ids = append(ids, gift.Id)
				probs = append(probs, float64(gift.Count))
			}
		}
		// 如果没有奖品，则返回"谢谢参与"
		if len(ids) == 0 {
			// CloseChannel() //关闭channel
			go CloseMQ()                               //关闭写mq的连接
			ctx.String(http.StatusOK, strconv.Itoa(0)) //0表示所有奖品已抽完
			return
		}
		// 抽奖算法，可手动选择哪种算法
		if lotteryAlgorithm == 1 {
			index = util.Lottery(probs) //抽中第index个奖品
		} else if lotteryAlgorithm == 2 {
			index = util.LotteryRandom(probs)
		} else if lotteryAlgorithm == 3 {
			// weights 是传入的每个奖品对应的权重值，可以联合其他业务，比如社交平台用户的访问次数，或用户购物次数等
			var weights []int
			index = util.LotteryWeightedRandom(probs, weights)
		}
		// 如果索引=-1，则表示抽奖失败，重试
		if index == -1 {
			continue
		}
		giftId := ids[index]
		err := database.ReduceInventory(giftId) // 先从redis上减库存
		if err != nil {
			util.LogRus.Warnf("奖品%d减库存失败", giftId)
			// 减库存失败，则重试
			continue
		} else {
			// 用户ID写死为1
			ProduceOrder(1, giftId) //把订单信息写入mq
			// 将mysql中的库存也减少
			ctx.String(http.StatusOK, strconv.Itoa(giftId)) //减库存成功后才给前端返回奖品ID
			return
		}
	}
	// 如果3次之后还失败，则返回“谢谢参与”
	ctx.String(http.StatusOK, strconv.Itoa(database.EMPTY_GIFT))
}

// Filter 拦截规则
func Filter(ctx *gin.Context) {

}
