package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetAccountPool(userId int64, currency string) (*models.AccountPool, error) {
	return mysql.SharedStore().GetAccountPool(userId, currency)
}

func GetAccountsPoolByUserId(userId int64) ([]*models.AccountPool, error) {
	return mysql.SharedStore().GetAccountsPoolByUserId(userId)
}
