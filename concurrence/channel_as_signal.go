package main

import (
	"fmt"
	"reflect"
	"time"
)

func main6() {
	ch := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("子协程结束")
		ch <- struct{}{}
	}()
	<-ch

	testEmptyStruct()
}

type A struct{}
type B struct{}

// 空结构体不占内存，获取它们的地址时go会返回统一的值
func testEmptyStruct() {
	a := A{}
	b := B{}
	fmt.Printf("%p   %p\n", &a, &b) //空结构体的地址是一样的
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	fmt.Printf("%d   %d\n", typeA.Size(), typeB.Size()) //空结构体的size都是0
}
