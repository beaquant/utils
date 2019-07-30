package wx

import (
	"testing"
	"time"
)

func TestWxPush_SendWxString(t *testing.T) {
	pusher := NewWxPush("https://sc.ftqq.com/", "SCU28704T55ea4ee1b39512b35eb63b36a24")
	pusher.SendWxString("test", "hello")
	time.Sleep(time.Second * 3)
}
