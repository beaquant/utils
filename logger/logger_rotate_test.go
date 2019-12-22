package logger

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestNewLoggerWithRotate(t *testing.T) {
	lr := NewLoggerWithRotate(".", "test_log", 1, time.Millisecond,
		&logrus.TextFormatter{ForceColors: true, TimestampFormat: SimpleDateTimeFormat, FullTimestamp: true}, logrus.DebugLevel, nil)
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
