package tasks

import (
	"time"
)

// YTL兑换BITE资金池任务
func StartConvertRelease() {
	ConvertRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			ConvertRelease()
		}
	}
}

func ConvertRelease() {
	_ = convertRelease()
}

func convertRelease() error {
	return nil
}
