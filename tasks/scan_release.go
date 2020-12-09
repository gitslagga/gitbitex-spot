package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
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
	// 查看当天是否释放过
	addressRelease, err := mysql.SharedStore().GetLastAddressRelease(models.AddressReleaseScan)
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
	sumFee, err := mysql.SharedStore().GetAccountScanSumFee()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ScanRelease GetAccountScanSumFee err: %v", err)
		return
	}
	if sumFee.LessThanOrEqual(decimal.Zero) {
		mylog.DataLogger.Info().Msgf("ScanRelease sumFee less than or equal decimal zero")
		return
	}

	count, err := mysql.SharedStore().CountAddressByGroupBite()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ScanRelease CountAddressByGroupBite err: %v", err)
		return
	}
	if count <= 0 {
		mylog.DataLogger.Info().Msgf("ScanRelease count less than or equal zero")
		return
	}

	addresses, err := mysql.SharedStore().GetAddressByGroupBite()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ScanRelease GetAddressByGroupBite err: %v", err)
		return
	}

	// 获取实际分红数量
	amount := sumFee.Mul(decimal.NewFromFloat(models.AccountScanReleaseRate)).Div(decimal.NewFromInt(int64(count)))
	for _, address := range addresses {
		err = scanRelease(address.Id, amount)
		if err != nil {
			mylog.DataLogger.Error().Msgf("ScanRelease scanRelease err: %v", err)
		}
	}
}

func scanRelease(userId int64, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyBite)
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
		Coin:   models.AccountCurrencyBite,
		Number: amount,
		Type:   models.AddressReleaseScan,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
