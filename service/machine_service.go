package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
)

func GetBuyMachine() ([]*models.Machine, error) {
	return mysql.SharedStore().GetBuyMachine()
}

func GetMachineById(machineId int64) (*models.Machine, error) {
	return mysql.SharedStore().GetMachineById(machineId)
}

func BuyMachine(address *models.Address, machine *models.Machine, currency string) error {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	var amount decimal.Decimal
	if currency == models.CURRENCY_YTL {
		rate, err := decimal.NewFromString(configs[15].Value)
		if err != nil {
			return err
		}
		amount = machine.Number.Mul(rate)
	} else if currency == models.CURRENCY_ENERGY {
		rate, err := decimal.NewFromString(configs[16].Value)
		if err != nil {
			return err
		}
		amount = machine.Number.Mul(rate)
	} else if currency == models.CURRENCY_USDT {
		amount = machine.Number
	} else {
		return errors.New("无效的币种|Invalid of currency")
	}

	return buyMachine(address, machine, currency, amount)
}

func buyMachine(address *models.Address, machine *models.Machine, currency string, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(address.Id, currency)
	if err != nil {
		return err
	}

	if addressAsset.Available.LessThan(amount) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	addressAsset.Available = addressAsset.Available.Sub(amount)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	err = db.AddMachineAddress(&models.MachineAddress{
		MachineId:   machine.Id,
		UserId:      address.Id,
		Number:      machine.Number.Add(machine.Number.Mul(machine.Profit)).Div(decimal.NewFromInt(int64(machine.Release))),
		TotalNumber: machine.Number.Add(machine.Number.Mul(machine.Profit)),
		Day:         machine.Release,
		TotalDay:    machine.Release,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
