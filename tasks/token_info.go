package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"time"
)

func StartTokenInfo() {
	t := time.NewTicker(60 * time.Minute)
	TokenInfo()

	for {
		select {
		case <-t.C:
			TokenInfo()
		}
	}
}

func TokenInfo() {
	config, err := mysql.SharedStore().GetAddressConfigByCoin(models.CurrencyWalletUsdt)
	if err != nil {
		mylog.DataLogger.Error().Msgf("TokenInfo GetAddressConfigByCoin err: %v", err)
		return
	}

	minDeposit = config.MinDeposit

	mylog.DataLogger.Info().Msgf("TokenInfo GetAddressConfigByCoin minDeposit %f", minDeposit)
}
