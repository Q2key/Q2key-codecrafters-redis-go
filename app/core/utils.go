package core

import (
	"encoding/json"
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"math/rand"
	"time"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func TraceObj(obj *contracts.Scheduler, label string, color string) {
	b, err := json.MarshalIndent((obj), "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\r\n%s%s:%s%s", color+label, time.Now().Local().UTC(), string(b), Reset)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
