package tasks

import (
	"time"
)

// 拼团节点资金池任务
func StartGroupRelease() {
	GroupRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			GroupRelease()
		}
	}
}

func GroupRelease() {
	_ = groupRelease()
}

func groupRelease() error {
	return nil
}
