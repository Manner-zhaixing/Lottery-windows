package main

import (
	"fmt"
	"gift/database"
	"gift/util"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

func init() {
	util.InitLog("log")
}

func TryLock(rc *redis.Client, name string, life time.Duration) bool {
	cmd := rc.SetNX(name, 1, life) //SetNX 如果key(name)不存在，则能成功写入，返回true
	if cmd.Err() != nil {
		return false //如果发生异常，认为上锁失败
	}
	return cmd.Val()
}

func ReleaseLock(rc *redis.Client, name string) {
	rc.Del(name)
}

func main12() {
	rc := database.GetRedisClient()
	const P = 100
	const lockName = "dqq"
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func(index int) {
			defer wg.Done()
			if TryLock(rc, lockName, 3*time.Minute) {
				fmt.Printf("协程%d 上锁成功\n", index)
			}
		}(i)
	}
	wg.Wait()
	ReleaseLock(rc, lockName)
}
