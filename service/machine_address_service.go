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

func MachineRelease() {
	t := time.NewTicker(6 * time.Hour)

	for {
		select {
		case <-t.C:
			machineAddress, err := mysql.SharedStore().GetMachineAddressUsedList()
			if len(machineAddress) > 0 && err == nil {
				for i := 0; i < len(machineAddress); i++ {
					//未结束的矿机才能挖矿
					if machineAddress[i].Day > 0 {
						machineLog, err := mysql.SharedStore().GetLastMachineLog(machineAddress[i].Id)
						if err != nil {
							mylog.DataLogger.Error().Msgf("MachineRelease GetLastMachineLog err: %v", err)
							continue
						}
						if machineLog != nil {
							t, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
							if err != nil {
								mylog.DataLogger.Error().Msgf("MachineRelease time parse err: %v", err)
								continue
							}

							if machineLog.UpdatedAt.After(t) && machineLog.UpdatedAt.Before(t.Add(24*time.Hour)) {
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

	addressAsset, err := db.GetAccountAssetForUpdate(machineAddress.UserId, models.CURRENCY_ENERGY)
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
