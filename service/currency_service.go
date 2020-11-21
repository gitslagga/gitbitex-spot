package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
)

func GetValidAddressConfig() ([]*models.AddressConfig, error) {
	return mysql.SharedStore().GetValidAddressConfig()
}

func GetAddressConfigByCoin(coin string) (*models.AddressConfig, error) {
	return mysql.SharedStore().GetAddressConfigByCoin(coin)
}

func UpdateAddressConfig(config *models.AddressConfig) error {
	return mysql.SharedStore().UpdateAddressConfig(config)
}

func AddAddressCollect(collect *models.AddressCollect) error {
	return mysql.SharedStore().AddAddressCollect(collect)
}

func GetAddressDepositsByUserId(userId, before, after, limit int64) ([]*models.AddressDeposit, error) {
	return mysql.SharedStore().GetAddressDepositsByUserId(userId, before, after, limit)
}

func AddAddressDeposit(deposit *models.AddressDeposit) error {
	return mysql.SharedStore().AddAddressDeposit(deposit)
}

func UpdateAddressDeposit(deposit *models.AddressDeposit) error {
	return mysql.SharedStore().UpdateAddressDeposit(deposit)
}

func GetAddressWithdrawsByUserId(userId, before, after, limit int64) ([]*models.AddressWithdraw, error) {
	return mysql.SharedStore().GetAddressWithdrawsByUserId(userId, before, after, limit)
}

func AddAddressWithdraw(withdraw *models.AddressWithdraw) error {
	return mysql.SharedStore().AddAddressWithdraw(withdraw)
}

func UpdateAddressWithdraw(withdraw *models.AddressWithdraw) error {
	return mysql.SharedStore().UpdateAddressWithdraw(withdraw)
}

func AddressWithdraw(address *models.Address, config *models.AddressConfig, toAddress string, numberF float64) error {
	number := decimal.NewFromFloat(numberF)
	if number.LessThan(config.MinWithdraw) {
		return errors.New("低于最小提币数量|Less than the minimum withdrawal amount")
	}

	actualNumber := number.Mul(decimal.New(1, 0).Add(config.WithdrawFee))

	return addressWithdraw(address, config, toAddress, number, actualNumber)
}

func addressWithdraw(address *models.Address, config *models.AddressConfig, toAddress string, number, actualNumber decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(address.Id, config.Coin)
	if err != nil {
		return err
	}

	if coinAsset.Available.LessThan(actualNumber) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	coinAsset.Available = coinAsset.Available.Sub(actualNumber)
	coinAsset.Hold = coinAsset.Hold.Add(actualNumber)
	err = db.UpdateAccountAsset(coinAsset)
	if err != nil {
		return err
	}

	err = db.AddAddressWithdraw(&models.AddressWithdraw{
		UserId:   address.Id,
		BlockNum: 0,
		TxId:     "",
		Coin:     config.Coin,
		Address:  toAddress,
		Value:    number,
		Actual:   actualNumber,
		Status:   models.CurrencyDepositUnConfirm,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
