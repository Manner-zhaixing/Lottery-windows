package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// 创建一个全局的 context 对象
var ctx = context.Background()

func main() {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 服务器地址
		Password: "123456",         // 设置 Redis 服务器密码
		DB:       0,                // 使用的数据库索引，默认是 0
	})

	// 测试连接
	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")

	// 设置一个键值对
	err = rdb.Set(ctx, "key", "ldddddddd", 0).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}

	// 获取刚刚设置的值
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatalf("Could not get key: %v", err)
	}
	fmt.Printf("key: %s\n", val)

}
