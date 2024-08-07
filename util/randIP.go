package util

import (
	"fmt"
	"gift/config"
	"math/rand"
	"strconv"
	"time"
)

// GenerateRandomIP 该文件产生随机IP地址
func GenerateRandomIP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d.%d.%d.%d",
		r.Intn(256),
		r.Intn(256),
		r.Intn(256),
		r.Intn(256),
	)
}

// GenerateRandomUserID 该文件产生随机用户ID
func GenerateRandomUserID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	loginUserID := strconv.Itoa(r.Intn(config.GenerateRandomMaxUserID))
	return loginUserID
}
