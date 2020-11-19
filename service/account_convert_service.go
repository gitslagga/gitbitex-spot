package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
	"time"
)

func GetAccountConvertByUserId(userId int64) ([]*models.AccountConvert, error) {
	return mysql.SharedStore().GetAccountConvertByUserId(userId)
}

func GetAccountConvertSumNumber() (float64, error) {
	sumNumber, err := mysql.SharedStore().GetAccountConvertSumNumber()
	if err != nil {
		return 0, err
	}

	number, _ := sumNumber.Float64()
	return number, nil
}

func AddAccountConvert(accountConvert *models.AccountConvert) error {
	return mysql.SharedStore().AddAccountConvert(accountConvert)
}

func AccountConvert(address *models.Address, convertType int, num float64) error {
	var err error
	switch convertType {
	case models.ConvertYtlToBite:
		err = AccountYtlConvertBite(address, num)
	case models.ConvertBiteToYtl:
		err = AccountBiteConvertYtl(address, num)
	}

	return err
}

func AccountYtlConvertBite(address *models.Address, num float64) error {
	number := decimal.NewFromFloat(num)
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	count, err := mysql.SharedStore().GetAccountConvertSumNumber()
	if err != nil {
		return err
	}

	total, err := decimal.NewFromString(configs[models.TotalYtlConvert].Value)
	if err != nil {
		return err
	}

	if count.Add(number).GreaterThan(total) {
		return errors.New("每日兑换数量超过限制|Daily convert quantity exceeds limit")
	}

	ytlRate, err := decimal.NewFromString(configs[models.RateYtlConvertUsdt].Value)
	if err != nil {
		return err
	}
	biteRate, err := decimal.NewFromString(configs[models.RateBiteConvertUsdt].Value)
	if err != nil {
		return err
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) || biteRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("兑换比例配置错误|Convert rate setting error")
	}

	price := ytlRate.Div(biteRate)
	amount := number.Div(price).Mul(address.ConvertFee.Add(decimal.New(1, 0)))

	err = accountYtlConvertBite(address, number, price, amount)
	if err != nil {
		return err
	}

	sumFee, err := mysql.SharedStore().GetAccountConvertSumFee()
	if err == nil {
		sumFeeFloat, _ := sumFee.Float64()
		currentTime := time.Now()
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

		_ = models.SharedRedis().SetAccountConvertSumFee(sumFeeFloat, endTime.Sub(currentTime))
	}

	return nil
}

func accountYtlConvertBite(address *models.Address, number, price, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	ytlAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CurrencyYtl)
	if err != nil {
		return err
	}
	if ytlAsset.Available.LessThan(amount) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	ytlAsset.Available = ytlAsset.Available.Sub(amount)
	err = db.UpdateAccountAsset(ytlAsset)
	if err != nil {
		return err
	}

	biteAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CurrencyBite)
	if err != nil {
		return err
	}
	biteAsset.Available = biteAsset.Available.Add(number)
	err = db.UpdateAccountAsset(biteAsset)
	if err != nil {
		return err
	}

	err = db.AddAccountConvert(&models.AccountConvert{
		UserId: address.Id,
		Number: number,
		Price:  price,
		Fee:    address.ConvertFee,
		Amount: amount,
		Type:   models.ConvertYtlToBite,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func AccountBiteConvertYtl(address *models.Address, num float64) error {
	number := decimal.NewFromFloat(num)
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	ytlRate, err := decimal.NewFromString(configs[models.RateYtlConvertUsdt].Value)
	if err != nil {
		return err
	}
	biteRate, err := decimal.NewFromString(configs[models.RateBiteConvertUsdt].Value)
	if err != nil {
		return err
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) || biteRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("兑换比例配置错误|Convert rate setting error")
	}

	price := biteRate.Div(ytlRate)
	amount := number.Div(price)

	err = accountBiteConvertYtl(address, number, price, amount)
	if err != nil {
		return err
	}

	return nil
}

func accountBiteConvertYtl(address *models.Address, number, price, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	biteAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CurrencyBite)
	if err != nil {
		return err
	}
	if biteAsset.Available.LessThan(amount) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	biteAsset.Available = biteAsset.Available.Sub(amount)
	err = db.UpdateAccountAsset(biteAsset)
	if err != nil {
		return err
	}

	ytlAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CurrencyYtl)
	if err != nil {
		return err
	}
	ytlAsset.Available = ytlAsset.Available.Add(number)
	err = db.UpdateAccountAsset(ytlAsset)
	if err != nil {
		return err
	}

	err = db.AddAccountConvert(&models.AccountConvert{
		UserId: address.Id,
		Number: number,
		Price:  price,
		Fee:    decimal.Zero,
		Amount: amount,
		Type:   models.ConvertBiteToYtl,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
