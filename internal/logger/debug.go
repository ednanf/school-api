package logger

import "log"

func DebugLogger(error error) {
	log.Printf("[DEBUG] %s\n", error)
}
