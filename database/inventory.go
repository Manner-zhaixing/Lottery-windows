package database

import (
	"fmt"
	"gift/util"
	"gorm.io/gorm"
	"strconv"
)

const (
	prefix = "gift_count_" //所有key设置统一的前缀，方便后续按前缀遍历key
)

// InitGiftInventory 从Mysql中读出所有奖品的初始库存，存入Redis。
// 如果同时有很多用户来参与抽奖活动，不能交发去Mysql里减库存，mysql扛不住这么高的并发量，Redis可以扛住
func InitGiftInventory() {
	giftCh := make(chan Gift, 100)
	// 将mysql中的数据读取并输入到giftCh管道中
	go GetAllGiftsV2(giftCh)
	// 获取redis客户端
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
		err := client.Set(prefix+strconv.Itoa(gift.Id), gift.Count, 0).Err() //0表示不设过期时间
		if err != nil {
			util.LogRus.Panicf("set gift %d:%s count to %d failed: %s", gift.Id, gift.Name, gift.Count, err)
		}
	}
}

// GetAllGiftInventory 获取所有奖品剩余的库存量Redis
func GetAllGiftInventory() []*Gift {
	client := GetRedisClient()
	keys, err := client.Keys(prefix + "*").Result() //根据前缀获取所有奖品的key
	if err != nil {
		util.LogRus.Errorf("iterate all keys by prefix %s failed: %s", prefix, err)
		return nil
	}
	gifts := make([]*Gift, 0, len(keys))
	for _, key := range keys { //根据奖品key获得奖品的库存count
		if id, err := strconv.Atoi(key[len(prefix):]); err == nil {
			count, err := client.Get(key).Int()
			if err == nil {
				gifts = append(gifts, &Gift{Id: id, Count: count})
			} else {
				util.LogRus.Errorf("invalid gift inventory %s", client.Get(key).String())
			}
		} else {
			util.LogRus.Errorf("invalid redis key %s", key)
		}
	}

	return gifts
}

// ReduceInventory 奖品对应的库存减1--redis
func ReduceInventory(GiftId int) error {
	client := GetRedisClient()
	key := prefix + strconv.Itoa(GiftId)
	// redis.Decr是原子操作
	n, err := client.Decr(key).Result()
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
// 后加的----begin
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

//---end
