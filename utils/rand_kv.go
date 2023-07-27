package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	randStr = rand.New(rand.NewSource(time.Now().Unix()))
	letters = []byte("afedfefgg12feawdawfdawf")
)

func GetTestKey(i int) []byte {
	return []byte(fmt.Sprintf("bitcask-go-key-%09d", i))
}

// RandomValue 生成随机value，用于测试
func RandomValue(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randStr.Intn(len(letters))]
	}
	return []byte("bitcask-go-value-" + string(b))
}
