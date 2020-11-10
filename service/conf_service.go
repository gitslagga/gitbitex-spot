package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetConfigs() ([]*models.Config, error) {
	return mysql.SharedStore().GetConfigs()
}

func GetConfigById(id int64) (*models.Config, error) {
	return mysql.SharedStore().GetConfigById(id)
}

func UpdateConfig(config *models.Config) error {
	return mysql.SharedStore().UpdateConfig(config)
}
