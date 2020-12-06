package service

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetAddressGroupByUserId(userId int64) ([]*models.AddressGroup, error) {
	return mysql.SharedStore().GetAddressGroupByUserId(userId)
}

func GetAddressGroupForUpdate(currency string) (*models.AddressGroup, error) {
	return mysql.SharedStore().GetAddressGroupForUpdate(currency)
}

func AddAddressGroup(group *models.AddressGroup) error {
	return mysql.SharedStore().AddAddressGroup(group)
}

func UpdateAddressGroup(group *models.AddressGroup) error {
	return mysql.SharedStore().UpdateAddressGroup(group)
}

func AddressGroup(address *models.Address, coin string) error {
	return nil
}

func AddressDelegate(address *models.Address, coin string) error {
	return nil
}

func AddressRelease(address *models.Address, coin string) error {
	return nil
}
