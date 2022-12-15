package logger

import (
	"fmt"
)

func Info(message interface{}) {
	fmt.Printf("[Info] %v\n", message)
}

func Warn(message interface{}) {
	fmt.Printf("[Warn] %v\n", message)
}

func Error(message interface{}) {
	fmt.Printf("[Error] %v\n", message)
}
