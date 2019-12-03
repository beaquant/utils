package logger

import (
	"testing"
	"time"
)

func TestNewLoggerWithRotate(t *testing.T) {
	lr := NewLoggerWithRotate(".", "test_log", 1, time.Millisecond, nil)
	cnt := 100
	for {
		lr.Println("hello")
		time.Sleep(time.Second)
		if cnt > 0 {
			cnt--
		} else {
			break
		}
	}
}
