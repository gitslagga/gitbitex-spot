package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetMachineLogByUserId(userId, before, after, limit int64) ([]*models.MachineLog, error) {
	return mysql.SharedStore().GetMachineLogByUserId(userId, before, after, limit)
}

func GetLastMachineLog(machineAddressId int64) (*models.MachineLog, error) {
	return mysql.SharedStore().GetLastMachineLog(machineAddressId)
}

func AddMachineLog(machineLog *models.MachineLog) error {
	return mysql.SharedStore().AddMachineLog(machineLog)
}
