package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetMachineLevel() ([]*models.MachineLevel, error) {
	return mysql.SharedStore().GetMachineLevel()
}

func GetMachineLevelById(machineLevelId int64) (*models.MachineLevel, error) {
	return mysql.SharedStore().GetMachineLevelById(machineLevelId)
}
