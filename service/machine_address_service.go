package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"time"
)

func GetMachineAddressByUserId(userId int64) ([]*models.MachineAddress, error) {
	return mysql.SharedStore().GetMachineAddressByUserId(userId)
}

func AddMachineAddress(machineAddress *models.MachineAddress) error {
	return mysql.SharedStore().AddMachineAddress(machineAddress)
}

func UpdateMachineAddress(machineAddress *models.MachineAddress) error {
	return mysql.SharedStore().UpdateMachineAddress(machineAddress)
}

func StartMachineRelease() {
	t := time.NewTicker(6 * time.Hour)
	MachineRelease()

	for {
		select {
		case <-t.C:
			MachineRelease()
		}
	}
}

func MachineRelease() {
	machineAddress, err := mysql.SharedStore().GetMachineAddressUsedList()
	if err == nil && len(machineAddress) > 0 {
		for i := 0; i < len(machineAddress); i++ {
			//获取矿机最后一条挖矿记录
			machineLog, err := mysql.SharedStore().GetLastMachineLog(machineAddress[i].Id)
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

			err = machineRelease(machineAddress[i])
			if err != nil {
				mylog.DataLogger.Error().Msgf("MachineRelease machineRelease err: %v", err)
			}
		}
	}
}

func machineRelease(machineAddress *models.MachineAddress) error {
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
		Number:           machineAddress.Number,
	})
	if err != nil {
		return err
	}

	addressAsset, err := db.GetAccountAssetForUpdate(machineAddress.UserId, models.CURRENCY_YTL)
	if err != nil {
		return err
	}

	addressAsset.Available = addressAsset.Available.Add(machineAddress.Number)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
