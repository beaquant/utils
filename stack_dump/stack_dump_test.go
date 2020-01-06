package stack_dump

import (
	"testing"
	"time"
)

func TestSetupStackTrap(t *testing.T) {
	SetupStackTrap()
	time.Sleep(100 * time.Second)
}
