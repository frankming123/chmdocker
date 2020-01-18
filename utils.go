package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// 生成一个sha256随机数
func gensha256() string {
	now := time.Now().UnixNano()
	sum := sha256.Sum256([]byte(string(now)))
	return fmt.Sprintf("%x",sum)
}
