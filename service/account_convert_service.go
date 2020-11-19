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

func AccountConvert(address *models.Address, num float64) error {
	number := decimal.NewFromFloat(num)
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	count, err := mysql.SharedStore().GetAccountConvertSumNumber()
	if err != nil {
		return err
	}

	total, err := decimal.NewFromString(configs[models.YtlConvertNumber].Value)
	if err != nil {
		return err
	}

	if count.Add(number).GreaterThan(total) {
		return errors.New("每日兑换数量超过限制|Daily convert quantity exceeds limit")
	}

	ytlRate, err := decimal.NewFromString(configs[models.YtlConvertUsdtRate].Value)
	if err != nil {
		return err
	}
	biteRate, err := decimal.NewFromString(configs[models.BiteConvertUsdtRate].Value)
	if err != nil {
		return err
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) || biteRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("兑换比例配置错误|Convert rate setting error")
	}

	price := ytlRate.Div(biteRate)
	amount := number.Div(price).Mul(address.ConvertFee.Add(decimal.New(1, 0)))

	err = accountConvert(address, number, price, amount)
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

func accountConvert(address *models.Address, number, price, amount decimal.Decimal) error {
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
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
