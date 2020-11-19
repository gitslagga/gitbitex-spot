package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
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
	case models.AccountCurrencyYtl:
		rate, err := decimal.NewFromString(configs[models.ConfigYtlConvertUsdt].Value)
		if err != nil {
			return err
		}
		if rate.LessThanOrEqual(decimal.Zero) {
			return errors.New("YTL兑换USDT价格错误|YTL convert USDT price error")
		}
		amount = machine.Number.Div(rate)
	case models.AccountCurrencyBite:
		rate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
		if err != nil {
			return err
		}
		if rate.LessThanOrEqual(decimal.Zero) {
			return errors.New("BITE兑换USDT价格错误|BITE convert USDT price error")
		}
		amount = machine.Number.Div(rate)
	case models.AccountCurrencyUsdt:
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
		IsBuy:       models.MachineBuy,
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
	parentAddressAsset, err := db.GetAccountAssetForUpdate(address.ParentId, models.AccountCurrencyYtl)
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

// 获取用户培养的达人数量，有效推荐数量(活跃度大于1)，总活跃度，小区活跃度
func getConditionMachineLevel(address *models.Address) (bool, error) {

	return true, nil
}

func StartMachineLevel(userId int64) {
	address, err := GetAddressById(userId)
	if err != nil {
		mylog.DataLogger.Error().Msgf("StartMachineLevel GetAddressById err: %v", err)
		return
	}

	condition, err := getConditionMachineLevel(address)
	if err != nil {
		mylog.DataLogger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
		return
	}
	if !condition {
		return
	}

	machineLevelStepOne(address)
}

func machineLevelStepOne(address *models.Address) {
	switch address.MachineLevelId {
	case models.MachineLevelZero:
		err := machineLevelStepTwo(address.Id, models.MachineLevelZero, models.MachineLevelOne)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelOne:
		err := machineLevelStepTwo(address.Id, models.MachineLevelOne, models.MachineLevelTwo)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelTwo:
		err := machineLevelStepTwo(address.Id, models.MachineLevelTwo, models.MachineLevelThree)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelThree:
		err := machineLevelStepTwo(address.Id, models.MachineLevelThree, models.MachineLevelFour)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelFour:
		err := machineLevelStepTwo(address.Id, models.MachineLevelFour, models.MachineLevelFive)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	}
}

// 进行分红，赠送矿机
func machineLevelStepTwo(userId, machineLevelId, countMachineLevelId int64) error {
	// 获取兑换手续费
	sumFee, err := models.SharedRedis().GetAccountConvertSumFee()
	if err != nil {
		return err
	}

	// 获取要升级的达人级别
	machineLevel, err := GetMachineLevelById(machineLevelId)
	if err != nil {
		return err
	}

	// 获取升级后的达人级别的数量
	count, err := CountAddressByMachineLevelId(countMachineLevelId)
	if err != nil {
		return err
	}

	// 获取实际分红数量
	var amount decimal.Decimal
	if count > 0 {
		amount = sumFee.Mul(machineLevel.GlobalFee).Div(decimal.NewFromInt(int64(count)))
	} else {
		amount = sumFee.Mul(machineLevel.GlobalFee)
	}

	// 获取赠送的矿机
	machine, err := mysql.SharedStore().GetMachineById(machineLevel.MachineId)
	if err != nil {
		return err
	}

	return machineLevelStepThree(userId, amount, machine)
}

func machineLevelStepThree(userId int64, amount decimal.Decimal, machine *models.Machine) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	ytlAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyYtl)
	if err != nil {
		return err
	}

	ytlAsset.Available = ytlAsset.Available.Add(amount)
	err = db.UpdateAccountAsset(ytlAsset)
	if err != nil {
		return err
	}

	err = db.AddMachineAddress(&models.MachineAddress{
		MachineId:   machine.Id,
		UserId:      userId,
		Number:      machine.Number.Add(machine.Number.Mul(machine.Profit)).Div(decimal.NewFromInt(int64(machine.Release))),
		TotalNumber: machine.Number.Add(machine.Number.Mul(machine.Profit)),
		Day:         machine.Release,
		TotalDay:    machine.Release,
		IsBuy:       models.MachineFree,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
