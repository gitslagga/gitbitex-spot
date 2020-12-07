package tasks

import (
	"time"
)

// 扫码支付资金池任务
func StartScanRelease() {
	ScanRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			ScanRelease()
		}
	}
}

func ScanRelease() {
	_ = scanRelease()
}

func scanRelease() error {
	return nil
}
