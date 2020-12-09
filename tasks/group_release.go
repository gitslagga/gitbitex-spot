package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
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
	// 查看当天是否释放过
	addressRelease, err := mysql.SharedStore().GetLastAddressRelease(models.AddressReleaseGroup)
	if err != nil {
		mylog.DataLogger.Info().Msgf("ConvertRelease GetLastAddressRelease err: %v", err)
		return
	}
	if addressRelease != nil {
		currentTime := time.Now()
		startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 00, 00, 00, 00, currentTime.Location())
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

		if addressRelease.CreatedAt.After(startTime) && addressRelease.CreatedAt.Before(endTime) {
			return
		}
	}

	// 获取扫描手续费
	usdtNum, err := mysql.SharedStore().GetAddressGroupSumNum(models.AccountGroupCurrencyUsdt)
	if err != nil {
		mylog.DataLogger.Error().Msgf("GroupRelease GetAddressGroupSumNum err: %v", err)
		return
	}
	biteNum, err := mysql.SharedStore().GetAddressGroupSumNum(models.AccountGroupCurrencyBite)
	if err != nil {
		mylog.DataLogger.Error().Msgf("GroupRelease GetAddressGroupSumNum err: %v", err)
		return
	}

	if usdtNum.GreaterThan(decimal.Zero) {
		count, err := mysql.SharedStore().CountAddressByGroupUsdt()
		if err != nil {
			mylog.DataLogger.Error().Msgf("GroupRelease CountAddressByGroupUsdt err: %v", err)
			return
		}
		if count <= 0 {
			mylog.DataLogger.Info().Msgf("GroupRelease count less than or equal zero")
			return
		}

		addresses, err := mysql.SharedStore().GetAddressByGroupUsdt()
		if err != nil {
			mylog.DataLogger.Error().Msgf("GroupRelease GetAddressByGroupUsdt err: %v", err)
			return
		}

		// 获取实际分红数量
		amount := biteNum.Mul(decimal.NewFromFloat(models.AccountGroupNodeRate)).Div(decimal.NewFromInt(int64(count)))
		for _, address := range addresses {
			err = groupRelease(address.Id, amount, models.AccountGroupCurrencyUsdt)
			if err != nil {
				mylog.DataLogger.Error().Msgf("GroupRelease groupRelease err: %v", err)
			}
		}
	}

	if biteNum.GreaterThan(decimal.Zero) {
		count, err := mysql.SharedStore().CountAddressByGroupBite()
		if err != nil {
			mylog.DataLogger.Error().Msgf("GroupRelease CountAddressByGroupBite err: %v", err)
			return
		}
		if count <= 0 {
			mylog.DataLogger.Info().Msgf("GroupRelease count less than or equal zero")
			return
		}

		addresses, err := mysql.SharedStore().GetAddressByGroupBite()
		if err != nil {
			mylog.DataLogger.Error().Msgf("GroupRelease GetAddressByGroupBite err: %v", err)
			return
		}

		// 获取实际分红数量
		amount := biteNum.Mul(decimal.NewFromFloat(models.AccountGroupNodeRate)).Div(decimal.NewFromInt(int64(count)))
		for _, address := range addresses {
			err = groupRelease(address.Id, amount, models.AccountGroupCurrencyBite)
			if err != nil {
				mylog.DataLogger.Error().Msgf("GroupRelease groupRelease err: %v", err)
			}
		}
	}
}

func groupRelease(userId int64, amount decimal.Decimal, coin string) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(userId, coin)
	if err != nil {
		return err
	}

	accountAsset.Available = accountAsset.Available.Add(amount)
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	err = db.AddAddressRelease(&models.AddressRelease{
		UserId: userId,
		Coin:   coin,
		Number: amount,
		Type:   models.AddressReleaseGroup,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
