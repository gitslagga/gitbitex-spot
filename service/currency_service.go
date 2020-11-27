package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/tasks"
	"github.com/gitslagga/gitbitex-spot/utils"
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

func GetAddressWithdrawsByOrderSN(orderSN string) (*models.AddressWithdraw, error) {
	return mysql.SharedStore().GetAddressWithdrawsByOrderSN(orderSN)
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
		OrderSN:  utils.GetOrderSN(),
		Value:    number,
		Actual:   actualNumber,
		Status:   models.CurrencyWithdrawReview,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

//backend
func BackendWithdraw(withdraw *models.AddressWithdraw, status int) error {
	if withdraw.Status != models.CurrencyWithdrawReview {
		return errors.New("订单已经处理|Order has been processed")
	}
	switch status {
	case models.CurrencyWithdrawSuccess:
		txId, err := sendTransactionWithoutFee(tasks.EthMainAddress, withdraw.Address, withdraw.Coin, withdraw.Value.String())
		if err != nil {
			mylog.Logger.Error().Msgf("AddressPassWithdraw PersonalSendTransactionToBlock err:%v", err)
			return err
		}

		withdraw.TxId = txId
		withdraw.Status = models.CurrencyWithdrawSuccess
		withdraw.BlockNum, err = tasks.EthBlockNumber()
		if err != nil {
			return err
		}

		err = backendWithdraw(withdraw)
		if err != nil {
			return err
		}
	case models.CurrencyWithdrawPassed:
		withdraw.Status = models.CurrencyWithdrawPassed
		err := backendWithdraw(withdraw)
		if err != nil {
			return err
		}
	case models.CurrencyWithdrawUnPass:
		withdraw.Status = models.CurrencyWithdrawUnPass
		err := backendUnPassWithdraw(withdraw)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendTransactionWithoutFee(fromAddress, toAddress, token, amount string) (string, error) {
	raw, err := tasks.PersonalSignTransaction(fromAddress, toAddress, token, amount)
	if err != nil {
		mylog.Logger.Error().Msgf("sendTransactionWithoutFee PersonalSignTransaction err:%v", err)
		return "", err
	}

	txId, err := tasks.EthSendRawTransaction(raw)
	if err != nil {
		mylog.Logger.Error().Msgf("sendTransactionWithoutFee EthSendRawTransaction err:%v", err)
		return "", err
	}

	return txId, nil
}

func backendWithdraw(withdraw *models.AddressWithdraw) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(withdraw.UserId, withdraw.Coin)
	if err != nil {
		return err
	}

	coinAsset.Hold = coinAsset.Hold.Sub(withdraw.Value)
	err = db.UpdateAccountAsset(coinAsset)
	if err != nil {
		return err
	}

	err = db.UpdateAddressWithdraw(withdraw)
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func backendUnPassWithdraw(withdraw *models.AddressWithdraw) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(withdraw.UserId, withdraw.Coin)
	if err != nil {
		return err
	}

	coinAsset.Available = coinAsset.Available.Add(withdraw.Value)
	coinAsset.Hold = coinAsset.Hold.Sub(withdraw.Value)
	err = db.UpdateAccountAsset(coinAsset)
	if err != nil {
		return err
	}

	err = db.UpdateAddressWithdraw(withdraw)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
