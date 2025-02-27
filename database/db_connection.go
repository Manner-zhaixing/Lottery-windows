package database

import (
	"context"
	"fmt"
	"gift/util"
	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var (
	blog_mysql      *gorm.DB
	blog_mysql_once sync.Once
	dblog           ormlog.Interface

	blog_redis      *redis.Client
	blog_redis_once sync.Once
)

func init() {
	dblog = ormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		ormlog.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel:      ormlog.Silent,
			Colorful:      false,
		},
	)
}

func createMysqlDB(dbname, host, user, pass string, port int) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbname)
	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dblog, PrepareStmt: true})
	if err != nil {
		util.LogRus.Panicf("mysql链接失败，panic: %s", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	util.LogRus.Infof("链接mysql成功 %s", dbname)
	return db
}

func GetGiftDBConnection() *gorm.DB {
	blog_mysql_once.Do(func() {
		if blog_mysql == nil {
			dbName := "gift"
			viper := util.CreateConfig("mysql")
			host := viper.GetString(dbName + ".host")
			port := viper.GetInt(dbName + ".port")
			user := viper.GetString(dbName + ".user")
			pass := viper.GetString(dbName + ".pass")
			blog_mysql = createMysqlDB(dbName, host, user, pass, port)
		}
	})

	return blog_mysql
}

func createRedisClient(address, pass string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: pass,
		DB:       db,
	})
	ctx := context.Background()
	if err := cli.Ping(ctx).Err(); err != nil {
		util.LogRus.Panicf("redis链接失败，panic: %s", err)
	} else {
		util.LogRus.Infof("connect to redis success")
	}
	return cli
}

func GetRedisClient() *redis.Client {
	blog_redis_once.Do(func() {
		if blog_redis == nil {
			viper := util.CreateConfig("redis")
			addr := viper.GetString("addr")
			pass := viper.GetString("password")
			db := viper.GetInt("db")
			blog_redis = createRedisClient(addr, pass, db)
		}
	})
	return blog_redis
}
