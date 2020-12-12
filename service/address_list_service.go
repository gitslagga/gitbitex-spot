package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
)

func GetAddressListByAddress(addressList string) (*models.AddressList, error) {
	return mysql.SharedStore().GetAddressListByAddress(addressList)
}

func GetAddressListById(userId int64) (*models.AddressList, error) {
	return mysql.SharedStore().GetAddressListById(userId)
}

func GetAddressListByUserId(userId int64) ([]*models.AddressList, error) {

	return mysql.SharedStore().GetAddressListByUserId(userId)
}

func CountAddressListByUserId(userId int64) (int, error) {
	return mysql.SharedStore().CountAddressListByUserId(userId)
}

func AddAddressList(addressList *models.AddressList) error {
	return mysql.SharedStore().AddAddressList(addressList)
}

func UpdateAddressList(addressList *models.AddressList) error {
	return mysql.SharedStore().UpdateAddressList(addressList)
}

func AddressListService(address *models.Address) ([]*models.AddressList, error) {
	addressList, err := GetAddressListByUserId(address.Id)
	if err != nil {
		return nil, err
	}

	return addressList, nil
}

func AddressAddList(userId int64, username, password, mnemonic string) error {
	count, err := CountAddressListByUserId(userId)
	if err != nil {
		return err
	}
	if count >= models.AccountMaxAddress {
		return errors.New("子地址添加数量超过限制|The number of sub addresses added exceeds the limit")
	}

	address, err := createAddressByMnemonic(mnemonic)
	if err != nil {
		return err
	}

	addressExists, err := mysql.SharedStore().GetAddressByAddress(address.Address)
	if err != nil {
		return err
	}
	addressListExists, err := mysql.SharedStore().GetAddressListByAddress(address.Address)
	if err != nil {
		return err
	}
	if addressExists != nil || addressListExists != nil {
		return errors.New("地址已存在|Address already exists")
	}

	addressList := &models.AddressList{
		UserId:     userId,
		Username:   username,
		Password:   password,
		Address:    address.Address,
		PublicKey:  address.PublicKey,
		PrivateKey: address.PrivateKey,
		Mnemonic:   address.Mnemonic,
	}
	return mysql.SharedStore().AddAddressList(addressList)
}

func DeleteAddressList(address string) error {
	addressList, err := GetAddressListByAddress(address)
	if err != nil {
		return err
	}
	if addressList == nil {
		return errors.New("子地址未找到|Sub address not found")
	}

	return mysql.SharedStore().DeleteAddressList(addressList.Id)
}

func AddressSwitchList(address *models.Address, addressList *models.AddressList) (*models.Address, error) {
	addressTemp := &models.Address{
		Username:   address.Username,
		Password:   address.Password,
		Address:    address.Address,
		PublicKey:  address.PublicKey,
		PrivateKey: address.PrivateKey,
		Mnemonic:   address.Mnemonic,
	}

	address.Username = addressList.Username
	address.Password = addressList.Password
	address.Address = addressList.Address
	address.PublicKey = addressList.PublicKey
	address.PrivateKey = addressList.PrivateKey
	address.Mnemonic = addressList.Mnemonic

	addressList.Username = addressTemp.Username
	addressList.Password = addressTemp.Password
	addressList.Address = addressTemp.Address
	addressList.PublicKey = addressTemp.PublicKey
	addressList.PrivateKey = addressTemp.PrivateKey
	addressList.Mnemonic = addressTemp.Mnemonic

	err := addressSwitchList(address, addressList)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func addressSwitchList(address *models.Address, addressList *models.AddressList) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	err = db.UpdateAddressList(addressList)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
