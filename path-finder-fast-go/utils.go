package main

import "fmt"

func Between(low, value, high int) bool {
	return low <= value && value < high
}

func DebugLog(format string, args ...interface{}) {
	fmt.Printf("[log] "+format+"\n", args...)
}
