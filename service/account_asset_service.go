package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetAccountAsset(userId int64, currency string) (*models.AccountAsset, error) {
	return mysql.SharedStore().GetAccountAsset(userId, currency)
}

func GetAccountsAssetByUserId(userId int64) ([]*models.AccountAsset, error) {
	return mysql.SharedStore().GetAccountsAssetByUserId(userId)
}
