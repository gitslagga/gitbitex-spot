package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

func GetBuyMachine() ([]*models.Machine, error) {
	return mysql.SharedStore().GetBuyMachine()
}

func GetMachineById(machineId int64) (*models.Machine, error) {
	return mysql.SharedStore().GetMachineById(machineId)
}

func BuyMachine(address *models.Address, machine *models.Machine, currency string) error {
	count, err := mysql.SharedStore().GetMachineAddressUsedCount(address.Id, machine.Id)
	if err != nil {
		return err
	}

	if count >= machine.BuyQuantity {
		return errors.New("可买数量受限|Available quantity limited")
	}

	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	var amount decimal.Decimal
	switch currency {
	case models.CURRENCY_YTL:
		rate, err := decimal.NewFromString(configs[15].Value)
		if err != nil {
			return err
		}
		if rate.LessThanOrEqual(decimal.Zero) {
			return errors.New("YTL兑换USDT价格错误|YTL convert USDT price error")
		}
		amount = machine.Number.Div(rate)
	case models.CURRENCY_BITE:
		rate, err := decimal.NewFromString(configs[17].Value)
		if err != nil {
			return err
		}
		if rate.LessThanOrEqual(decimal.Zero) {
			return errors.New("BITE兑换USDT价格错误|BITE convert USDT price error")
		}
		amount = machine.Number.Div(rate)
	case models.CURRENCY_USDT:
		amount = machine.Number
	default:
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
		IsBuy:       1,
	})
	if err != nil {
		return err
	}

	//增加活跃度
	address.ActiveNum = address.ActiveNum + machine.Active
	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	//增加上级直推奖励
	parentIds := strings.Split(address.ParentIds, "-")
	parentId, err := strconv.ParseInt(parentIds[len(parentIds)-1], 10, 64)
	if err != nil {
		return err
	}

	parentAddressAsset, err := db.GetAccountAssetForUpdate(parentId, models.CURRENCY_YTL)
	if err != nil {
		return err
	}

	parentAddressAsset.Available = parentAddressAsset.Available.Add(machine.Number.Mul(machine.Invite))
	err = db.UpdateAccountAsset(parentAddressAsset)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
