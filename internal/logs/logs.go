package logs

import (
	"log"
	"time"
)

func GetNow() string {
	now := time.Now().Format("02 Jan 2006 15:04:00")
	return now
}

func LogEvent(err error) {
	now := GetNow()
	log.Printf("[ %s ] context: %v", now, err)
}
