package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetMachineAddressByUserId(userId, before, after, limit int64) ([]*models.MachineAddress, error) {
	return mysql.SharedStore().GetMachineAddressByUserId(userId, before, after, limit)
}

func AddMachineAddress(machineAddress *models.MachineAddress) error {
	return mysql.SharedStore().AddMachineAddress(machineAddress)
}

func UpdateMachineAddress(machineAddress *models.MachineAddress) error {
	return mysql.SharedStore().UpdateMachineAddress(machineAddress)
}
