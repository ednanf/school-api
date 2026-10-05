package logger

import "fmt"

func DebugLogger(error error) {
	fmt.Printf("[DEBUG] %s\n", error)
}
