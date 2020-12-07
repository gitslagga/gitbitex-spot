package service

import (
	"errors"
	"fmt"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/utils"
	"github.com/shopspring/decimal"
)

func ExecuteBill(userId int64, currency string) error {
	tx, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 锁定用户资金记录
	account, err := tx.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	// 资金记录不存在，创建一条，并再次执行加锁
	if account == nil {
		err = tx.AddAccount(&models.Account{
			UserId:    userId,
			Currency:  currency,
			Available: decimal.Zero,
		})
		if err != nil {
			return err
		}
		account, err = tx.GetAccountForUpdate(userId, currency)
		if err != nil {
			return err
		}
	}

	// 获取所有未入账的bill
	bills, err := tx.GetUnsettledBillsByUserId(userId, currency)
	if err != nil {
		return err
	}
	if len(bills) == 0 {
		return nil
	}

	for _, bill := range bills {
		account.Available = account.Available.Add(bill.Available)
		account.Hold = account.Hold.Add(bill.Hold)

		bill.Settled = true

		err = tx.UpdateBill(bill)
		if err != nil {
			return err
		}
	}

	err = tx.UpdateAccount(account)
	if err != nil {
		return err
	}

	err = tx.CommitTx()
	if err != nil {
		return err
	}

	return nil
}

func HoldBalance(db models.Store, userId int64, currency string, size decimal.Decimal, billType models.BillType) error {
	if size.LessThanOrEqual(decimal.Zero) {
		return errors.New("size less than 0")
	}

	enough, err := HasEnoughBalance(userId, currency, size)
	if err != nil {
		return err
	}
	if !enough {
		return errors.New(fmt.Sprintf("no enough %v : request=%v", currency, size))
	}

	account, err := db.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("no enough")
	}

	account.Available = account.Available.Sub(size)
	account.Hold = account.Hold.Add(size)

	bill := &models.Bill{
		UserId:    userId,
		Currency:  currency,
		Available: size.Neg(),
		Hold:      size,
		Type:      billType,
		Settled:   true,
		Notes:     "",
	}
	err = db.AddBills([]*models.Bill{bill})
	if err != nil {
		return err
	}

	err = db.UpdateAccount(account)
	if err != nil {
		return err
	}

	return nil
}

func HasEnoughBalance(userId int64, currency string, size decimal.Decimal) (bool, error) {
	account, err := GetAccount(userId, currency)
	if err != nil {
		return false, err
	}
	if account == nil {
		return false, nil
	}
	return account.Available.GreaterThanOrEqual(size), nil
}

func GetAccount(userId int64, currency string) (*models.Account, error) {
	return mysql.SharedStore().GetAccount(userId, currency)
}

func GetAccountsByUserId(userId int64) ([]*models.Account, error) {
	return mysql.SharedStore().GetAccountsByUserId(userId)
}

func AddDelayBill(store models.Store, userId int64, currency string, available, hold decimal.Decimal, billType models.BillType, notes string) (*models.Bill, error) {
	bill := &models.Bill{
		UserId:    userId,
		Currency:  currency,
		Available: available,
		Hold:      hold,
		Type:      billType,
		Settled:   false,
		Notes:     notes,
	}
	err := store.AddBills([]*models.Bill{bill})
	return bill, err
}

func GetUnsettledBills() ([]*models.Bill, error) {
	return mysql.SharedStore().GetUnsettledBills()
}

func AccountAddress(userId int64) (map[string]interface{}, error) {
	accountSpot, err := GetAccountsByUserId(userId)
	if err != nil {
		return nil, err
	}

	accountAsset, err := GetAccountsAssetByUserId(userId)
	if err != nil {
		return nil, err
	}

	accountPool, err := GetAccountsPoolByUserId(userId)
	if err != nil {
		return nil, err
	}

	accountShop, err := GetAccountsShopByUserId(userId)
	if err != nil {
		return nil, err
	}

	configs, err := GetConfigs()
	if err != nil {
		return nil, err
	}

	ytlRate, err := decimal.NewFromString(configs[models.ConfigYtlConvertUsdt].Value)
	if err != nil {
		return nil, err
	}
	biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
	if err != nil {
		return nil, err
	}
	usdtRate, err := decimal.NewFromString(configs[models.ConfigUsdtConvertCny].Value)
	if err != nil {
		return nil, err
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) || biteRate.LessThanOrEqual(decimal.Zero) || usdtRate.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("兑换比例配置错误|Convert rate setting error")
	}

	var calculateUsdt decimal.Decimal
	accountAssetMap := make([]map[string]interface{}, len(accountAsset))
	accountPoolMap := make([]map[string]interface{}, len(accountPool))
	accountSpotMap := make([]map[string]interface{}, len(accountSpot))
	accountShopMap := make([]map[string]interface{}, len(accountShop))
	for k, v := range accountAsset {
		accountAssetMap[k] = utils.StructToMapViaJson(v)
		switch v.Currency {
		case models.AccountCurrencyYtl:
			calculateUsdt = calculateUsdt.Add(ytlRate.Mul(v.Available.Add(v.Hold)))
			accountAssetMap[k]["Calculate"] = usdtRate.Mul(ytlRate.Mul(v.Available.Add(v.Hold)))
		case models.AccountCurrencyBite:
			calculateUsdt = calculateUsdt.Add(biteRate.Mul(v.Available.Add(v.Hold)))
			accountAssetMap[k]["Calculate"] = usdtRate.Mul(biteRate.Mul(v.Available.Add(v.Hold)))
		case models.AccountCurrencyUsdt:
			calculateUsdt = calculateUsdt.Add(v.Available.Add(v.Hold))
			accountAssetMap[k]["Calculate"] = usdtRate.Mul(v.Available.Add(v.Hold))
		}
	}

	for k, v := range accountPool {
		accountPoolMap[k] = utils.StructToMapViaJson(v)
		switch v.Currency {
		case models.AccountCurrencyBite:
			calculateUsdt = calculateUsdt.Add(biteRate.Mul(v.Available.Add(v.Hold)))
			accountPoolMap[k]["Calculate"] = usdtRate.Mul(biteRate.Mul(v.Available.Add(v.Hold)))
		}
	}

	for k, v := range accountSpot {
		accountSpotMap[k] = utils.StructToMapViaJson(v)
		switch v.Currency {
		case models.AccountCurrencyBite:
			calculateUsdt = calculateUsdt.Add(biteRate.Mul(v.Available.Add(v.Hold)))
			accountSpotMap[k]["Calculate"] = usdtRate.Mul(biteRate.Mul(v.Available.Add(v.Hold)))
		case models.AccountCurrencyUsdt:
			calculateUsdt = calculateUsdt.Add(v.Available.Add(v.Hold))
			accountSpotMap[k]["Calculate"] = usdtRate.Mul(v.Available.Add(v.Hold))
		}
	}

	for k, v := range accountShop {
		accountShopMap[k] = utils.StructToMapViaJson(v)
		switch v.Currency {
		case models.AccountCurrencyBite:
			calculateUsdt = calculateUsdt.Add(biteRate.Mul(v.Available.Add(v.Hold)))
			accountShopMap[k]["Calculate"] = usdtRate.Mul(biteRate.Mul(v.Available.Add(v.Hold)))
		case models.AccountCurrencyUsdt:
			calculateUsdt = calculateUsdt.Add(v.Available.Add(v.Hold))
			accountShopMap[k]["Calculate"] = usdtRate.Mul(v.Available.Add(v.Hold))
		}
	}

	calculateCny := calculateUsdt.Mul(usdtRate)

	return map[string]interface{}{
		"accountSpot":   accountSpotMap,
		"accountAsset":  accountAssetMap,
		"accountPool":   accountPoolMap,
		"accountShop":   accountShopMap,
		"calculateUsdt": calculateUsdt,
		"calculateCny":  calculateCny,
	}, nil
}
