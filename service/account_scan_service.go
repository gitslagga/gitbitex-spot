package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
	"time"
)

func GetAccountScanByUserId(userId int64) ([]*models.AccountScan, error) {
	return mysql.SharedStore().GetAccountScanByUserId(userId)
}

func GetAccountScanSumNumber(userId int64) (decimal.Decimal, error) {
	return mysql.SharedStore().GetAccountScanSumNumber(userId)
}

func GetAccountScanSumFee() (decimal.Decimal, error) {
	return mysql.SharedStore().GetAccountScanSumFee()
}

func AddAccountScan(accountScan *models.AccountScan) error {
	return mysql.SharedStore().AddAccountScan(accountScan)
}

func AccountScan(userId int64, url string, numberF float64) error {
	sumNumber, err := GetAccountScanSumNumber(userId)
	if err != nil {
		return err
	}
	if sumNumber.Add(decimal.NewFromFloat(numberF)).GreaterThan(decimal.NewFromInt(models.AccountScanDayPayment)) {
		return errors.New("每日支付额度超过限制|Daily payment exceeds the limit")
	}

	configs, err := GetConfigs()
	if err != nil {
		return err
	}

	usdtRate, err := decimal.NewFromString(configs[models.ConfigUsdtConvertCny].Value)
	if err != nil {
		return err
	}
	biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
	if err != nil {
		return err
	}
	if usdtRate.LessThanOrEqual(decimal.Zero) || biteRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("兑换比例配置错误|Convert rate setting error")
	}

	number := decimal.NewFromFloat(numberF)
	fee := decimal.NewFromFloat(0.05)
	actualNumber := decimal.NewFromFloat(numberF * (1 + 0.05))

	rate := usdtRate.Div(biteRate)
	amount := actualNumber.Mul(rate)

	err = accountScan(userId, url, number, fee, actualNumber, rate, amount)
	if err != nil {
		return err
	}

	sumFee, err := mysql.SharedStore().GetAccountScanSumFee()
	if err == nil {
		currentTime := time.Now()
		endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

		_ = models.SharedRedis().SetAccountScanSumFee(sumFee, endTime.Sub(currentTime))
	}

	return nil
}

func accountScan(userId int64, url string, number, fee, actualNumber, rate, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountScanCurrency)
	if err != nil {
		return err
	}

	if addressAsset.Available.LessThan(amount) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	addressAsset.Available = addressAsset.Available.Sub(amount)
	addressAsset.Hold = addressAsset.Hold.Add(amount)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	err = db.AddAccountScan(&models.AccountScan{
		UserId:       userId,
		Currency:     models.AccountScanCurrency,
		Url:          url,
		Number:       number,
		Fee:          fee,
		ActualNumber: actualNumber,
		Rate:         rate,
		Amount:       amount,
		Status:       models.AccountScanUnPayment,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
