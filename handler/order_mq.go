package handler

import (
	"context"
	"gift/database"
	"gift/util"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

var (
	writer    *kafka.Writer
	writeWg   sync.WaitGroup
	closeOnce sync.Once
)

func InitMQ() {
	viper := util.CreateConfig("kafka")
	writer = &kafka.Writer{
		Addr:                   kafka.TCP(viper.GetString("brokers")),
		Topic:                  viper.GetString("topic"),
		AllowAutoTopicCreation: true, //Topic不存在时自动创建
	}
	util.LogRus.Info("create writer to mq")
}

// ProduceOrder 把订单放入mq
func ProduceOrder(UserId, GiftId int) {
	order := database.Order{UserId: UserId, GiftId: GiftId}
	writeWg.Add(1)
	go func() { //写MQ太慢，异步执行
		defer writeWg.Done()
		json, _ := sonic.Marshal(order)
		if err := writer.WriteMessages(context.Background(), kafka.Message{Value: json}); err != nil {
			util.LogRus.Errorf("write kafka failed: %s", err)
		}
	}()
}

// CloseMQ 关闭mq连接。CloseMQ()可以被反复调用
func CloseMQ() {
	closeOnce.Do(func() { //确保只执行一次
		writeWg.Wait() //由于写mq是异步执行的，这里要等所有的写操作全部完成，才能关闭writer
		writer.Close()
		util.LogRus.Info("stop write mq")
	})
}
