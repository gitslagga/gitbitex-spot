package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetMachineConfigs() ([]*models.MachineConfig, error) {
	return mysql.SharedStore().GetMachineConfigs()
}

func GetMachineConfigById(id int64) (*models.MachineConfig, error) {
	return mysql.SharedStore().GetMachineConfigById(id)
}

func UpdateMachineConfig(config *models.MachineConfig) error {
	return mysql.SharedStore().UpdateMachineConfig(config)
}
