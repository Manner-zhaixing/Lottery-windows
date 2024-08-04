package handler

import (
	"gift/config"
	"gift/database"
	"gift/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	// 抽奖之前的验证逻辑
	// 1.验证登录用户

	// 2.用户抽奖的分布式锁
	// 根据用户id获取分布式锁，防止用户短时间大量点击抽奖按钮
	loginUserID := "111"
	flagGetDisLock := util.DisLock(loginUserID)
	if !flagGetDisLock {
		// 获取失败，返回
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 102,
			"msg":  "获取分布式锁失败,请勿重复点击",
			"data": -1,
		})
		return
	} else {
		// 获取分布式锁成功，defer关闭
		defer util.DisUnLock(loginUserID)
	}

	// 3.验证用户今日抽奖次数
	UsesSumsDay := database.IncrementUserSumsDay(loginUserID)
	if UsesSumsDay > config.UserSumsDayLimitMax {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 103,
			"msg":  "今日该用户抽奖次数已用完",
			"data": -1,
		})
		return
	}
	// 4.验证ip今日的抽奖次数
	clientIP := ctx.ClientIP()
	IPSumsDay := database.IncrementIpSumsDay(clientIP)
	if IPSumsDay > config.IPSumsDayLimitMax {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 104,
			"msg":  "今日该ip抽奖次数已用完",
			"data": -1,
		})
		return
	}
	// 5.验证ip黑名单
	flagInBanIP := database.CheckBanIPs(clientIP)
	if flagInBanIP {
		// 在黑名单中，拦截
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 105,
			"msg":  "该ip已被拉黑",
			"data": -1,
		})
		return
	}
	// 6.验证用户黑名单
	userIdInt, _ := strconv.Atoi(loginUserID)
	flagInBanUser := database.CheckBanUsers(userIdInt)
	if flagInBanUser {
		// 在黑名单中，拦截
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 105,
			"msg":  "该用户已被拉黑",
			"data": -1,
		})
		return
	}
	// 8.匹配奖品是否中奖，即抽奖逻辑
	giftId, giftInformation := lotteryLogic()
	if giftId == -1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 106,
			"msg":  "抽奖失败",
			"data": -1,
		})
		return
	} else {
		// 抽奖成功
		err := database.ReduceInventory(giftId) // 先从redis上减库存
		if err != nil {
			util.LogRus.Warnf("奖品%d减库存失败", giftId)
			// 减库存失败，则重试
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 107,
				"msg":  "减库存失败，抽奖失败",
				"data": -1,
			})
			return
		} else {
			// 11.记录中奖订单
			// 用户ID写死,把订单信息写入mq
			userID, _ := strconv.Atoi(loginUserID)
			ProduceOrder(userID, giftId, giftInformation)
			//减库存成功后才给前端返回奖品ID
			// 12.返回抽奖结果
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "抽奖成功",
				"data": giftId,
			})
			return
		}
	}
}

func lotteryLogic() (int, *database.Gift) {
	// 使用的抽奖算法类别 [1-库存作为权重,2-随机,3-手动赋值权重,可以结合其他业务]
	var index int
	// 抽奖算法中的奖品权重
	var weights []int
	// gifts := make([]*Gift, 0, len(keys))
	// Gift{Id: id, Count: count}
	gifts := database.GetAllGiftInventory() //获取所有奖品剩余的库存量Redis
	ids := make([]int, 0, len(gifts))
	probs := make([]float64, 0, len(gifts))
	giftsInformation := make([]*database.Gift, 0, len(gifts))
	for _, gift := range gifts {
		if gift.Count > 0 { //先确保redis返回的库存量大小0，因为抽奖算法Lottery不支持抽中概率为0的奖品
			ids = append(ids, gift.Id)
			probs = append(probs, float64(gift.Count))
			giftsInformation = append(giftsInformation, gift)
		}
	}
	// 如果没有奖品，则返回"谢谢参与"
	if len(ids) == 0 {
		// CloseChannel() //关闭channel
		go CloseMQ() //关闭写mq的连接
		return 0, nil
	}
	// 抽奖算法，可手动选择哪种算法
	if config.LotteryAlgorithm == 1 {
		index = util.Lottery(probs) //抽中第index个奖品
	} else if config.LotteryAlgorithm == 2 {
		index = util.LotteryRandom(probs)
	} else if config.LotteryAlgorithm == 3 {
		// weights 是传入的每个奖品对应的权重值，可以联合其他业务，比如社交平台用户的访问次数，或用户购物次数等
		index = util.LotteryWeightedRandom(probs, weights)
	}
	// 如果索引=-1，则表示抽奖失败
	if index == -1 {
		return -1, nil
	}
	giftId := ids[index]
	return giftId, giftsInformation[index]
}
