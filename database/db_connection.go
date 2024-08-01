package database

import (
	"fmt"
	"gift/util"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	ormlog "gorm.io/gorm/logger"
)

// 建立数据库连接。代码讲解参见《双Token博客系统》

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
		util.LogRus.Panicf("connect to mysql use dsn %s failed: %s", dsn, err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	util.LogRus.Infof("connect to mysql db %s", dbname)
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

func createRedisClient(address, passwd string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: passwd,
		DB:       db,
	})
	if err := cli.Ping().Err(); err != nil {
		util.LogRus.Panicf("connect to redis %d failed %v", db, err)
	} else {
		util.LogRus.Infof("connect to redis %d", db)
	}
	return cli
}

func GetRedisClient() *redis.Client {
	blog_redis_once.Do(func() {
		if blog_redis == nil {
			viper := util.CreateConfig("redis")
			addr := viper.GetString("addr")
			pass := viper.GetString("pass")
			db := viper.GetInt("db")
			blog_redis = createRedisClient(addr, pass, db)
		}
	})
	return blog_redis
}
