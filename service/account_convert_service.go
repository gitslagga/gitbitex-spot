package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
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

	total, err := decimal.NewFromString(configs[2].Value)
	if err != nil {
		return err
	}

	if count.Add(number).GreaterThan(total) {
		return errors.New("每日兑换数量超过限制|Daily convert quantity exceeds limit")
	}

	energyRate, err := decimal.NewFromString(configs[16].Value)
	if err != nil {
		return err
	}
	bitcRate, err := decimal.NewFromString(configs[17].Value)
	if err != nil {
		return err
	}
	if bitcRate.LessThanOrEqual(decimal.Zero) || energyRate.Div(bitcRate).LessThanOrEqual(decimal.Zero) {
		return errors.New("兑换比例配置错误|Convert rate setting error")
	}

	price := energyRate.Div(bitcRate)
	amount := number.Div(price).Mul(address.ConvertFee.Add(decimal.New(1, 0)))

	err = accountConvert(address, number, price, amount)
	if err != nil {
		return err
	}

	return nil
}

func accountConvert(address *models.Address, number, price, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	energyAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CURRENCY_ENERGY)
	if err != nil {
		return err
	}
	if energyAsset.Available.LessThan(amount) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	energyAsset.Available = energyAsset.Available.Sub(amount)
	err = db.UpdateAccountAsset(energyAsset)
	if err != nil {
		return err
	}

	bitcAsset, err := db.GetAccountAssetForUpdate(address.Id, models.CURRENCY_BITC)
	if err != nil {
		return err
	}
	bitcAsset.Available = bitcAsset.Available.Add(number)
	err = db.UpdateAccountAsset(bitcAsset)
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
