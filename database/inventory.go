package database

import (
	"context"
	"fmt"
	"gift/config"
	"gift/util"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var (
	prefix                = config.GiftCountPrefix //所有key设置统一的前缀，方便后续按前缀遍历key
	giftGTypeMinMaxPrefix = config.GiftGTypeMinMaxPrefix
)

// InitGiftInventory 从Mysql中读出所有奖品的初始库存，存入Redis。
// 如果同时有很多用户来参与抽奖活动，不能交发去Mysql里减库存，mysql扛不住这么高的并发量，Redis可以扛住
func InitGiftInventory() {
	giftCh := make(chan Gift, 100)
	// 将mysql中的数据读取并输入到giftCh管道中
	go GetAllGiftsV2(giftCh)
	// 获取redis客户端
	ctx := context.Background()
	client := GetRedisClient()
	for {
		gift, ok := <-giftCh
		if !ok { //channel已经消费完了
			break
		}
		if gift.Count <= 0 {
			// util.LogRus.Warnf("gift %d:%s count is %d", gift.Id, gift.Name, gift.Count)
			continue //没有库存的商品不参与抽奖
		}
		// 将奖品数据存入到redis,不设过期时间
		err := client.Set(ctx, prefix+strconv.Itoa(gift.Id), gift.Count, 0).Err() //0表示不设过期时间
		if err != nil {
			util.LogRus.Panicf("set gift %d:%s count to %d failed: %s", gift.Id, gift.Name, gift.Count, err)
		}
		// 将奖品的GType，minWeight,maxWeight拼接之后存入到redis
		giftGTypeMinMax := strconv.Itoa(gift.GType) + "_" + strconv.Itoa(gift.MinWeight) + "_" + strconv.Itoa(gift.MaxWeight)
		err = client.Set(ctx, giftGTypeMinMaxPrefix+strconv.Itoa(gift.Id), giftGTypeMinMax, 0).Err() //0表示不设过期时间
		if err != nil {
			util.LogRus.Panicf("set gift %d:%s count to %d failed: %s", gift.Id, gift.Name, gift.Count, err)
		}
	}
}

// GetAllGiftInventory 获取所有奖品剩余的库存量Redis
func GetAllGiftInventory() []*Gift {
	ctx := context.Background()
	client := GetRedisClient()
	keys, err := client.Keys(ctx, prefix+"*").Result() //根据前缀获取所有奖品的key
	if err != nil {
		util.LogRus.Errorf("iterate all keys by prefix %s failed: %s", prefix, err)
		return nil
	}
	if err != nil {
		util.LogRus.Errorf("iterate all keys by prefix %s failed: %s", giftGTypeMinMaxPrefix, err)
		return nil
	}
	gifts := make([]*Gift, 0, len(keys))
	for _, key := range keys { //根据奖品key获得奖品的库存count
		id, _ := strconv.Atoi(key[len(prefix):])
		count, errCount := client.Get(ctx, key).Int()
		others, errOthers := client.Get(ctx, giftGTypeMinMaxPrefix+strconv.Itoa(id)).Result()
		if errCount != nil && errOthers != nil {
			othersSlice := strings.Split(others, "_")
			if len(othersSlice) != 3 {
				util.LogRus.Errorf("[lottert read redis failed]")
			} else {
				temp := &Gift{
					Id:        id,
					Count:     count,
					GType:     util.StrToInt(othersSlice[0]),
					MinWeight: util.StrToInt(othersSlice[1]),
					MaxWeight: util.StrToInt(othersSlice[2]),
				}
				gifts = append(gifts, temp)
			}
		} else {
			util.LogRus.Errorf("[lottert read redis failed]")
		}
		//if id, err := strconv.Atoi(key[len(prefix):]); err == nil {
		//	count, err := client.Get(ctx, key).Int()
		//	if err == nil {
		//
		//	} else {
		//		util.LogRus.Errorf("invalid gift inventory %s", client.Get(ctx, key).String())
		//	}
		//} else {
		//	util.LogRus.Errorf("invalid redis key %s", key)
		//}
	}

	return gifts
}

// ReduceInventory 奖品对应的库存减1--redis
func ReduceInventory(GiftId int) error {
	client := GetRedisClient()
	ctx := context.Background()
	key := prefix + strconv.Itoa(GiftId)
	// redis.Decr是原子操作
	n, err := client.Decr(ctx, key).Result()
	if err != nil {
		util.LogRus.Errorf("decr key %s failed: %s", key, err)
		return err
	} else {
		if n < 0 {
			util.LogRus.Errorf("%d已无库存，减1失败", GiftId)
			return fmt.Errorf("%d已无库存，减1失败", GiftId)
		}
		return nil
	}
}

// ReduceInventoryMysql 奖品对应的库存减1--mysql
func ReduceInventoryMysql(GiftId int) error {
	mysqlClient := GetGiftDBConnection()
	err := mysqlClient.Model(&Gift{Id: GiftId}).Updates(map[string]interface{}{"Count": gorm.Expr("Count - ?", 1)})
	// err := mysqlClient.Where("id = ?", string(GiftId)).Find(&gifts)
	if err != nil {
		// util.LogRus.Errorf("update gift inventory %d failed: %s", GiftId, err)
		return fmt.Errorf("update gift inventory %d failed: %s", GiftId, err)
	}
	return nil
}
