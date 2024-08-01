package main

import (
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conn *gorm.DB
)

func InitConn() {
	var err error
	//建立连接
	conn, err = gorm.Open(mysql.Open("tester:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

type User struct {
	Id int
}

func (User) TableName() string {
	return "user"
}

func Query() {
	var user User
	if err := conn.Select("id").First(&user).Error; err == nil {
		fmt.Println(user.Id)
	} else {
		fmt.Println(err.Error())
	}
}

func main() {
	InitConn()
	Query()

	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			Query()
		}()
	}
	wg.Wait()
}

// go run ./concurrent_connection/client/gorm
