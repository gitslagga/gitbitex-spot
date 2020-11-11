package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetAccountShop(userId int64, currency string) (*models.AccountShop, error) {
	return mysql.SharedStore().GetAccountShop(userId, currency)
}

func GetAccountsShopByUserId(userId int64) ([]*models.AccountShop, error) {
	return mysql.SharedStore().GetAccountsShopByUserId(userId)
}
