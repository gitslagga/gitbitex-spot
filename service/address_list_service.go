package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/utils"
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

func AddressListService(address *models.Address) ([]map[string]interface{}, error) {
	addressList, err := GetAddressListByUserId(address.Id)
	if err != nil {
		return nil, err
	}

	addressListMap := make([]map[string]interface{}, len(addressList)+1)
	addressListMap[0] = utils.StructToMapViaJson(address)

	for k, _ := range addressList {
		addressListMap[k+1] = utils.StructToMapViaJson(addressList)
	}

	return addressListMap, nil
}

func AddressAddList(userId int64, mnemonic, privateKey, password string) error {
	count, err := CountAddressListByUserId(userId)
	if err != nil {
		return err
	}
	if count >= models.AccountMaxAddress {
		return errors.New("子地址添加数量超过限制|The number of sub addresses added exceeds the limit")
	}

	var address *models.Address
	if mnemonic != "" {
		address, err = createAddressByMnemonic(mnemonic)
	} else {
		address, err = createAddressByPrivateKey(privateKey)
	}
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
		Username:   models.AccountDefaultName,
		Password:   password,
		Address:    address.Address,
		PublicKey:  address.PublicKey,
		PrivateKey: privateKey,
		Mnemonic:   address.Mnemonic,
	}
	return mysql.SharedStore().AddAddressList(addressList)
}

func DeleteAddressList(address string) error {
	addressList, err := GetAddressListByAddress(address)
	if err != nil {
		return err
	}

	return mysql.SharedStore().DeleteAddressList(addressList)
}

func AddressSwitchList(address *models.Address, addressList *models.AddressList) error {
	addressTemp := *address

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

	return addressSwitchList(address, addressList)
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
