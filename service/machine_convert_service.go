package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
	"time"
)

func GetMachineConvertByUserId(userId int64) ([]*models.MachineConvert, error) {
	return mysql.SharedStore().GetMachineConvertByUserId(userId)
}

func GetMachineConvertSumNumber() (float64, error) {
	sumNumber, err := mysql.SharedStore().GetMachineConvertSumNumber()
	if err != nil {
		return 0, err
	}

	number, _ := sumNumber.Float64()
	return number, nil
}

func AddMachineConvert(machineConvert *models.MachineConvert) error {
	return mysql.SharedStore().AddMachineConvert(machineConvert)
}

func MachineConvert(address *models.Address, convertType int, num float64) error {
	var err error
	switch convertType {
	case models.ConvertYtlToBite:
		err = MachineYtlConvertBite(address, num)
	case models.ConvertBiteToYtl:
		err = MachineBiteConvertYtl(address, num)
	}

	return err
}

func MachineYtlConvertBite(address *models.Address, num float64) error {
	number := decimal.NewFromFloat(num)
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	count, err := mysql.SharedStore().GetMachineConvertSumNumber()
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

	err = machineYtlConvertBite(address, number, price, amount)
	if err != nil {
		return err
	}

	sumFee, err := mysql.SharedStore().GetMachineConvertSumFee()
	if err == nil {
		sumFeeFloat, _ := sumFee.Float64()
		currentTime := time.Now()
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

		_ = models.SharedRedis().SetMachineConvertSumFee(sumFeeFloat, endTime.Sub(currentTime))
	}

	return nil
}

func machineYtlConvertBite(address *models.Address, number, price, amount decimal.Decimal) error {
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

	err = db.AddMachineConvert(&models.MachineConvert{
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

func MachineBiteConvertYtl(address *models.Address, num float64) error {
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

	err = machineBiteConvertYtl(address, number, price, amount)
	if err != nil {
		return err
	}

	return nil
}

func machineBiteConvertYtl(address *models.Address, number, price, amount decimal.Decimal) error {
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

	err = db.AddMachineConvert(&models.MachineConvert{
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
