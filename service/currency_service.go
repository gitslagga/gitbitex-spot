package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
)

func GetValidCurrencies() ([]*models.Currency, error) {
	return mysql.SharedStore().GetValidCurrencies()
}

func GetCurrencyByCoin(coin string) (*models.Currency, error) {
	return mysql.SharedStore().GetCurrencyByCoin(coin)
}

func UpdateCurrency(currency *models.Currency) error {
	return mysql.SharedStore().UpdateCurrency(currency)
}

func AddCurrencyCollect(currencyCollect *models.CurrencyCollect) error {
	return mysql.SharedStore().AddCurrencyCollect(currencyCollect)
}

func GetCurrencyDepositsByUserId(userId int64) ([]*models.CurrencyDeposit, error) {
	return mysql.SharedStore().GetCurrencyDepositsByUserId(userId)
}

func AddCurrencyDeposit(currencyDeposit *models.CurrencyDeposit) error {
	return mysql.SharedStore().AddCurrencyDeposit(currencyDeposit)
}

func UpdateCurrencyDeposit(currencyDeposit *models.CurrencyDeposit) error {
	return mysql.SharedStore().UpdateCurrencyDeposit(currencyDeposit)
}

func GetCurrencyWithdrawsByUserId(userId int64) ([]*models.CurrencyWithdraw, error) {
	return mysql.SharedStore().GetCurrencyWithdrawsByUserId(userId)
}

func AddCurrencyWithdraw(currencyWithdraw *models.CurrencyWithdraw) error {
	return mysql.SharedStore().AddCurrencyWithdraw(currencyWithdraw)
}

func UpdateCurrencyWithdraw(currencyWithdraw *models.CurrencyWithdraw) error {
	return mysql.SharedStore().UpdateCurrencyWithdraw(currencyWithdraw)
}

func CurrencyWithdraw(address *models.Address, currency *models.Currency, toAddress string, numberF float64) error {
	number := decimal.NewFromFloat(numberF)
	if number.LessThan(currency.MinWithdraw) {
		return errors.New("低于最小提币数量|Less than the minimum withdrawal amount")
	}

	actualNumber := number.Mul(decimal.New(1, 0).Add(currency.WithdrawFee))

	return currencyWithdraw(address, currency, toAddress, number, actualNumber)
}

func currencyWithdraw(address *models.Address, currency *models.Currency, toAddress string, number, actualNumber decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(address.Id, currency.Coin)
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

	err = db.AddCurrencyWithdraw(&models.CurrencyWithdraw{
		UserId:   address.Id,
		BlockNum: 0,
		TxId:     "",
		Coin:     currency.Coin,
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
