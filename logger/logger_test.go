package logger

import (
	"testing"
	"time"
)

func TestInitLog(t *testing.T) {
	logger := NewLogger("1.txt")
	//logger.WithField("nihao", "good").Infoln("hello")
	logger.Infoln("bad")
	time.Sleep(5 * time.Second)
}
