package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
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
	// 查看当天是否释放过
	addressRelease, err := mysql.SharedStore().GetLastAddressRelease(models.AddressReleaseConvert)
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

	// 获取兑换手续费
	sumFee, err := mysql.SharedStore().GetMachineConvertSumFee()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ConvertRelease GetMachineConvertSumFee err: %v", err)
		return
	}
	if sumFee.LessThanOrEqual(decimal.Zero) {
		mylog.DataLogger.Info().Msgf("ConvertRelease sumFee less than or equal decimal zero")
		return
	}

	machineLevel, err := mysql.SharedStore().GetMachineLevel()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ConvertRelease GetMachineLevel err: %v", err)
		return
	}

	for _, val := range machineLevel {
		// 获取升级后的达人级别的数量
		count, err := mysql.SharedStore().CountAddressByMachineLevelId(val.Id)
		if err != nil {
			mylog.DataLogger.Error().Msgf("ConvertRelease CountAddressByMachineLevelId err: %v", err)
			return
		}
		if count <= 0 {
			mylog.DataLogger.Info().Msgf("ConvertRelease count less than or equal zero")
			continue
		}

		addresses, err := mysql.SharedStore().GetAddressByMachineLevelId(val.Id)
		if err != nil {
			mylog.DataLogger.Error().Msgf("ConvertRelease GetAddressByMachineLevelId err: %v", err)
			return
		}

		// 获取实际分红数量
		amount := sumFee.Mul(val.GlobalFee).Div(decimal.NewFromInt(int64(count)))
		for _, address := range addresses {
			err = convertRelease(address.Id, amount)
			if err != nil {
				mylog.DataLogger.Error().Msgf("ConvertRelease convertRelease err: %v", err)
			}
		}
	}
}

func convertRelease(userId int64, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyYtl)
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
		Coin:   models.AccountCurrencyYtl,
		Number: amount,
		Type:   models.AddressReleaseConvert,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
