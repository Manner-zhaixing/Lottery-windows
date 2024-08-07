package handler

import (
	"gift/config"
	"gift/database"
	"gift/service"
	"gift/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAllGifts 获取所有奖品信息，用于初始化轮盘
func GetAllGifts(ctx *gin.Context) {
	//原始代码：从mysql获取初始轮盘的奖品信息
	util.LogRus.Infof("[gift-GetAllGifts.func] 获取所有奖品信息")
	gifts := database.GetAllGiftsV1()
	//改进：从redis获取初始轮盘的奖品信息
	//gifts := database.GetAllGiftInventory()
	if len(gifts) == 0 {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
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
		return
	}
}

// Lottery 抽奖
func Lottery(ctx *gin.Context) {
	// 产生随机userID和IP地址
	loginUserID := util.GenerateRandomUserID()
	clientIP := util.GenerateRandomIP()
	// 抽奖之前的验证逻辑
	// 1.验证登录用户

	// 2.用户抽奖的分布式锁
	// 根据用户id获取分布式锁，防止用户短时间大量点击抽奖按钮
	//loginUserID := "112"
	flagGetDisLock := service.DisLock(loginUserID)
	if !flagGetDisLock {
		// 获取失败，返回
		util.LogRus.Warnf("[gift-lottery.func]获取分布式锁失败")
		ctx.JSON(http.StatusOK, gin.H{
			"code": 102,
			"msg":  "获取分布式锁失败,谢谢参与",
			"data": -1,
		})
		return
	} else {
		// 获取分布式锁成功，defer关闭
		defer service.DisUnLock(loginUserID)
	}

	// 3.验证用户今日抽奖次数
	UsesSumsDay := database.IncrementUserSumsDay(loginUserID)
	if UsesSumsDay > config.UserSumsDayLimitMax {
		util.LogRus.Warnf("[gift-lottery.func]今日该用户抽奖次数已用完")
		ctx.JSON(http.StatusOK, gin.H{
			"code": 103,
			"msg":  "今日该用户抽奖次数已用完",
			"data": -1,
		})
		return
	}
	// 4.验证ip今日的抽奖次数
	IPSumsDay := database.IncrementIpSumsDay(clientIP)
	if IPSumsDay > config.IPSumsDayLimitMax {
		util.LogRus.Warnf("[gift-lottery.func] 今日该ip抽奖次数已用完")
		ctx.JSON(http.StatusOK, gin.H{
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
	//// 6.验证用户黑名单
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

	// 8.匹配奖品是否中奖，即抽奖逻辑;返回值[-1 代表无库存，取消抽奖;1为谢谢参与;其余id为正常发奖]
	// 查看保底次数
	var guaranteeFlag = false
	guaranteeSumsNow := database.GuaranteeSumIncrement(loginUserID)
	if guaranteeSumsNow >= int64(config.GuaranteeSum) {
		//抽奖次数+1之后到了指定的保底次数,返回true，保底
		guaranteeFlag = true
	}
	giftId, giftInformation := lotteryLogic(BlackListBool, guaranteeFlag)
	if giftId == -1 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 106,
			"msg":  "无库存，取消抽奖",
			"data": -1,
		})
		return
	} else if giftId == 1 {
		// 抽到谢谢参与
		// 保底次数增加
		database.GuaranteeSumIncrement(loginUserID)
		util.LogRus.Infof("用户:%s,谢谢参与", loginUserID)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 106,
			"msg":  "谢谢参与",
			"data": 1,
		})
		return
	} else {
		// 抽奖成功
		err := database.ReduceInventory(giftId) // 先从redis上减库存
		if err != nil && giftId != 1 {
			util.LogRus.Warnf("奖品%d减库存失败", giftId)
			// 减库存失败，则重试
			util.LogRus.Infof("用户:%s,谢谢参与,redis减库存失败", loginUserID)
			ctx.JSON(http.StatusOK, gin.H{
				"code": 107,
				"msg":  "减库存失败，抽奖失败",
				"data": -1,
			})
			// 抽奖失败，保底次数增加
			database.GuaranteeSumIncrement(loginUserID)
			return
		} else {
			// 11.记录中奖订单
			// 用户ID写死
			userID, _ := strconv.Atoi(loginUserID)
			// 将id和user拉入黑名单
			mysqlClient := database.GetGiftDBConnection()
			_, err = database.CreateBanIP(mysqlClient, clientIP, config.DeadTime)
			if err != nil {
				util.LogRus.Warnf("mysql插入banIP失败,err:%s", err)
			}
			// 中了大奖才拉入黑名单，黑名单中的人不许抽大奖
			if giftInformation.GType == 0 {
				// 存在的话验证deadTime，过了黑名单的时间就更新，没过的话，不更新，不允许抽大奖，不存在黑名单就新建
				_, err = database.UpdateOrInsertBanUser(mysqlClient, userID, config.DeadTime)
				if err != nil {
					util.LogRus.Warnf("mysql插入banUser失败,err:%s", err)
				}
			}
			// 保底次数重置
			database.GuaranteeSumReset(loginUserID)
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
func lotteryLogic(BlackListBool bool, guaranteeFlag bool) (int, *database.Gift) {
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

	// 查看所有奖品的库存是否都为0
	var (
		sum             = 0
		giftID          int
		giftInformation *database.Gift
	)
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
	if !guaranteeFlag {
		// 非保底，正常抽奖
		giftID, giftInformation = database.LotteryWeightedRandom(gifts)
		if BlackListBool && giftInformation.GType == 0 {
			//黑名单抽到大奖，取消抽奖，改为谢谢参与
			return 1, nil
		}
	} else {
		// 保底，获取有库存且GType为0的奖品，进行抽奖
		// tempGifts 存储有库存且GType为0的奖品
		tempGifts := make([]*database.Gift, 0)
		for k, v := range gifts {
			if v.Count > 0 && v.GType == 0 {
				tempGifts = append(tempGifts, gifts[k])
			}
		}
		giftID, giftInformation = database.LotteryWeightedRandom(gifts)
	}

	// giftID如果为1，代表谢谢参与
	return giftID, giftInformation
}
