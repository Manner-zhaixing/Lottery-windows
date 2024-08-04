package main

import (
	"context"
	"fmt"
	"gift/config"
	"gift/database"
	"gift/util"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

var reader *kafka.Reader

func Init() {
	util.InitLog("mq_consumer_log")
	viper := util.CreateConfig("kafka")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(viper.GetString("brokers"), ","),
		Topic:          viper.GetString("topic"),
		StartOffset:    kafka.LastOffset,  //之前MQ里的老数据不再接收了
		GroupID:        "serialize_order", //注意：如果不指定GroupID，则只能消费到1个partition里的数据。我的kafka配的是生成2个partition
		CommitInterval: 1 * time.Second,   //每隔多长时间自动commit一次offset
	})
	util.LogRus.Info("create reader to mq")
}

// ConsumeOrder 从mq里取出订单，把订单写入Mysql
func ConsumeOrder() {
	for {
		if message, err := reader.ReadMessage(context.Background()); err != nil {
			fmt.Printf("read message from mq failed: %v", err)
			break
		} else {
			var order database.Order
			if err = sonic.Unmarshal(message.Value, &order); err == nil {
				util.LogRus.Debugf("message partition %d", message.Partition)
				// 将订单写入mysql
				database.CreateOrder(order)
				// 将mysql中的库存减一
				err = database.ReduceInventoryMysql(order.GiftId)
				if err != nil {
					util.LogRus.Errorf("[consumer] reduce inventory mysql failed: %v", err)
				}
				// 将该用户拉入黑名单，一定时间不能抽奖
				mysqlConn := database.GetGiftDBConnection()
				_, err = database.CreateBanUser(mysqlConn, order.UserId, int(config.BanUsersExpireTime))
				if err != nil {
					util.LogRus.Errorf("[consumer] user to banUsers failed")
				}
				// 将该ip拉入黑名单，一定时间不能抽奖
				//_, err = database.CreateBanIP(mysqlConn, clientIP, int(config.BanIPsExpireTime))
			} else {
				util.LogRus.Errorf("order info is invalid json format: %s", string(message.Value))
			}
		}
	}
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号。Ctrl+C对应SIGINT信号
	sig := <-c                                        //阻塞，直到信号的到来
	reader.Close()
	util.LogRus.Infof("receive signal %s, exit", sig.String())
	os.Exit(0) //进程退出

}

func main() {
	Init()
	go listenSignal()
	ConsumeOrder()
}

// go run ./mq_consumer/
