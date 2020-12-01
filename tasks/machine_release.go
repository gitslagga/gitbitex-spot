package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"time"
)

func StartMachineRelease() {
	MachineRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			MachineRelease()
		}
	}
}

func MachineRelease() {
	//获取YTL兑换USDT费率
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		mylog.DataLogger.Error().Msgf("MachineRelease GetConfigs err: %v", err)
		return
	}
	ytlRate, err := decimal.NewFromString(configs[models.ConfigYtlConvertUsdt].Value)
	if err != nil {
		mylog.DataLogger.Error().Msgf("MachineRelease ytlRate err: %v", err)
		return
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) {
		mylog.DataLogger.Error().Msgf("MachineRelease YTL convert USDT price error")
		return
	}

	machineAddress, err := mysql.SharedStore().GetMachineAddressUsedList()
	if err == nil {
		for _, v := range machineAddress {
			//获取矿机最后一条挖矿记录
			machineLog, err := mysql.SharedStore().GetLastMachineLog(v.Id)
			if err != nil {
				mylog.DataLogger.Error().Msgf("MachineRelease GetLastMachineLog err: %v", err)
				continue
			}
			if machineLog != nil {
				currentTime := time.Now()
				startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 00, 00, 00, 00, currentTime.Location())
				endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

				if machineLog.CreatedAt.After(startTime) && machineLog.CreatedAt.Before(endTime) {
					continue
				}
			}

			//获取实际应得的YTL数量
			number := v.Number.Div(ytlRate)

			err = machineRelease(v, number)
			if err != nil {
				mylog.DataLogger.Error().Msgf("MachineRelease machineRelease err: %v", err)
			}
		}
	}
}

func machineRelease(machineAddress *models.MachineAddress, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	machineAddress.Day--
	err = db.UpdateMachineAddress(machineAddress)
	if err != nil {
		return err
	}

	err = db.AddMachineLog(&models.MachineLog{
		UserId:           machineAddress.UserId,
		MachineId:        machineAddress.MachineId,
		MachineAddressId: machineAddress.Id,
		Number:           number,
	})
	if err != nil {
		return err
	}

	addressAsset, err := db.GetAccountAssetForUpdate(machineAddress.UserId, models.AccountCurrencyYtl)
	if err != nil {
		return err
	}

	addressAsset.Available = addressAsset.Available.Add(number)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
