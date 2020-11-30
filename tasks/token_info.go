package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"time"
)

func StartTokenInfo() {
	TokenInfo()

	t := time.NewTicker(60 * time.Minute)
	for {
		select {
		case <-t.C:
			TokenInfo()
		}
	}
}

func TokenInfo() {
	config, err := mysql.SharedStore().GetAddressConfigByCoin(UsdtName)
	if err != nil {
		mylog.DataLogger.Error().Msgf("TokenInfo GetAddressConfigByCoin err: %v", err)
		return
	}

	ethColdAddress2 = config.CollectAddress
}
