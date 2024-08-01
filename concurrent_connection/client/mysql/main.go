package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	conn *sql.DB
)

func InitConn() {
	var err error
	//建立连接
	conn, err = sql.Open("mysql", "tester:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")
	if err != nil {
		panic(err)
	}
}

func Query() {
	rows, _ := conn.Query("select uid from user limit 1")
	defer rows.Close()
	for rows.Next() {
		var uid int
		rows.Scan(&uid)
		fmt.Println(uid)
	}
}

func main() {
	InitConn()
	defer conn.Close()
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

// go run ./concurrent_connection/client/mysql
