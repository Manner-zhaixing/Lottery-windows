package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

const PRODUC_NUM = 3

var buffer chan string = make(chan string, 100)
var pc_sync = make(chan struct{}, PRODUC_NUM)
var all_over = make(chan struct{})

// 一个生产者，负责读一个文件，把每行的内容放入buffer
func producer(filename string) {
	fin, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fin.Close()

	reader := bufio.NewReader(fin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					buffer <- (line + "\n")
				}
				break
			} else {
				fmt.Println(err)
			}
		} else {
			buffer <- line
		}
	}
	<-pc_sync
}

func consumer(filename string) {
	fout, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)

	for {
		if len(buffer) == 0 { //1.生产者都结束了,消费者把buffer里的内容都消费完了。2.生产者还没结束，但是消费者把buffer里的内容都消费完了
			if len(pc_sync) == 0 { //所有生产者都结束了
				break
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			line := <-buffer
			writer.WriteString(line)
		}
	}
	writer.Flush()
	all_over <- struct{}{}
}

func main7() {
	for i := 0; i < PRODUC_NUM; i++ {
		pc_sync <- struct{}{}
	}

	go producer("data/1.txt")
	go producer("data/2.txt")
	go producer("data/3.txt")
	go consumer("data/big.txt")

	<-all_over
}
