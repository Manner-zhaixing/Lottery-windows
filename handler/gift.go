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
	//原始代码：从mysql获取初始轮盘的奖品信息
	gifts := database.GetAllGiftsV1()
	//改进：从redis获取初始轮盘的奖品信息
	//gifts := database.GetAllGiftInventory()
	if len(gifts) == 0 {
		ctx.JSON(http.StatusInternalServerError, nil)
	} else {
		//抹掉敏感信息
		for _, gift := range gifts {
			gift.Count = -1
			gift.Price = -1
			gift.GType = -1
			gift.MinWeight = -1
			gift.MaxWeight = -1
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
	var BlackListBool = false
	flagInBanIP := database.CheckBanIPs(clientIP)
	if flagInBanIP {
		// 在黑名单，限制该ip无法抽取大奖
		BlackListBool = true
		// 在黑名单中，拦截
		//ctx.JSON(http.StatusInternalServerError, gin.H{
		//	"code": 105,
		//	"msg":  "该ip已被拉黑",
		//	"data": -1,
		//})
	}
	// 6.验证用户黑名单
	userIdInt, _ := strconv.Atoi(loginUserID)
	if !BlackListBool {
		// 如果ip不在黑名单，验证用户黑名单
		flagInBanUser := database.CheckBanUsers(userIdInt)
		if flagInBanUser {
			//在用户黑名单中，该用户短时间不能抽取大奖
			BlackListBool = true
			// 在黑名单中，拦截
			//ctx.JSON(http.StatusInternalServerError, gin.H{
			//	"code": 105,
			//	"msg":  "该用户已被拉黑",
			//	"data": -1,
			//})
		}
	}

	// 8.匹配奖品是否中奖，即抽奖逻辑;返回值[-1代表无库存，取消抽奖;1为谢谢参与;其余id为正常发奖]
	giftId, giftInformation := lotteryLogic(BlackListBool)
	if giftId == -1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 106,
			"msg":  "无库存，取消抽奖",
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
			// 用户ID写死
			userID, _ := strconv.Atoi(loginUserID)
			// 将id和user拉入黑名单
			mysqlClient := database.GetGiftDBConnection()
			_, err := database.CreateBanIP(mysqlClient, clientIP, config.DeadTime)
			if err != nil {
				util.LogRus.Warnf("mysql插入banIP失败")
			}
			_, err = database.CreateBanUser(mysqlClient, userID, config.DeadTime)
			if err != nil {
				util.LogRus.Warnf("mysql插入banUser失败")
			}
			// 把订单信息写入mq
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

// lotteryLogic 抽奖逻辑
func lotteryLogic(BlackListBool bool) (int, *database.Gift) {
	// 使用的抽奖算法类别 [1-库存作为权重,2-随机,3-手动赋值权重,可以结合其他业务]
	// gifts := make([]*Gift, 0, len(keys))
	// Gift
	//{
	//  Id        int
	//	Count     int
	//	GType     int
	//	MinWeight int
	//	MaxWeight int
	//}
	gifts := database.GetAllGiftInventory() //获取所有奖品剩余的库存量Redis
	//ids := make([]int, 0, len(gifts))
	//probs := make([]float64, 0, len(gifts))

	// 查看所有奖品的库存是否都为0
	var sum int = 0
	for _, v := range gifts {
		if v.Count <= 0 {
			sum++
		}
	}
	if sum == len(gifts) {
		// 所有奖品都无库存
		go CloseMQ()
		return -1, nil
	}

	//// 抽奖算法，可手动选择哪种算法
	//if config.LotteryAlgorithm == 1 {
	//	giftID = util.Lottery(probs) //抽中第index个奖品
	//} else if config.LotteryAlgorithm == 2 {
	//	giftID = util.LotteryRandom(probs)
	//} else if config.LotteryAlgorithm == 3 {
	//	// weights 是传入的每个奖品对应的权重值，可以联合其他业务，比如社交平台用户的访问次数，或用户购物次数等
	//	giftID, giftInformation := util.LotteryWeightedRandom(gifts)
	//}
	giftID, giftInformation := util.LotteryWeightedRandom(gifts)
	if BlackListBool && giftInformation.GType == 0 {
		//黑名单抽到大奖，取消抽奖，改为谢谢参与
		return 1, nil
	}
	// giftID如果为1，代表谢谢参与
	return giftID, giftInformation
}
