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

// 判断升级所需要的用户培养达人数量，有效推荐数量(活跃度大于1)，总活跃度，小区活跃度(扣除活跃度最大的哪一枝)
func getConditionMachineLevel(address *models.Address, machineLevel *models.MachineLevel) (bool, error) {
	sonAddress, err := GetAddressByParentId(address.Id)
	if err != nil {
		return false, err
	}

	var trainNumberOne int
	var trainNumberTwo int
	var trainNumberThree int

	var inviteNum int
	var maxActiveNum int
	var totalActiveNum int
	for i := 0; i < len(sonAddress); i++ {
		if sonAddress[i].MachineLevelId == models.MachineLevelOne {
			trainNumberOne++
		}
		if sonAddress[i].MachineLevelId == models.MachineLevelTwo {
			trainNumberTwo++
		}
		if sonAddress[i].MachineLevelId == models.MachineLevelThree {
			trainNumberThree++
		}

		if sonAddress[i].ActiveNum > models.MachineValidInvite {
			inviteNum++
		}

		activeNum, err := getTotalActiveNumber(sonAddress[i].Id)
		if err != nil {
			return false, err
		}
		if activeNum > maxActiveNum {
			maxActiveNum = activeNum
		}
		totalActiveNum += activeNum
	}

	if inviteNum < machineLevel.InviteNum {
		return false, nil
	}

	switch address.MachineLevelId {
	case models.MachineLevelZero:
	case models.MachineLevelOne:
	case models.MachineLevelTwo:
		if trainNumberOne < machineLevel.TrainNum {
			return false, nil
		}
	case models.MachineLevelThree:
		if trainNumberTwo < machineLevel.TrainNum {
			return false, nil
		}
	case models.MachineLevelFour:
		if trainNumberThree < machineLevel.TrainNum {
			return false, nil
		}
	}

	if totalActiveNum < machineLevel.TotalActive {
		return false, nil
	}
	if totalActiveNum-maxActiveNum < machineLevel.CommonActive {
		return false, nil
	}

	return true, nil
}

func getTotalActiveNumber(userId int64) (int, error) {
	sonAddress, err := GetAddressByParentId(userId)
	if err != nil {
		return 0, err
	}

	var totalNum int
	var activeNum int
	for i := 0; i < len(sonAddress); i++ {
		activeNum, err = getTotalActiveNumber(sonAddress[i].Id)
		if err != nil {
			return 0, err
		}

		totalNum += activeNum
	}

	return totalNum, nil
}

func StartMachineLevel(userId int64) {
	address, err := GetAddressById(userId)
	if err != nil {
		mylog.DataLogger.Error().Msgf("StartMachineLevel GetAddressById err: %v", err)
		return
	}

	machineLevel, err := GetMachineLevel()
	if err != nil {
		mylog.DataLogger.Error().Msgf("StartMachineLevel GetMachineLevel err: %v", err)
		return
	}

	machineLevelStepOne(address, machineLevel)
}

func machineLevelStepOne(address *models.Address, machineLevel []*models.MachineLevel) {
	switch address.MachineLevelId {
	case models.MachineLevelZero:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelZero])
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelZero], address, models.MachineLevelOne)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelOne:
		err := machineLevelStepTwo(machineLevel[models.MachineLevelOne], address, models.MachineLevelTwo)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelTwo:
		err := machineLevelStepTwo(machineLevel[models.MachineLevelTwo], address, models.MachineLevelThree)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelThree:
		err := machineLevelStepTwo(machineLevel[models.MachineLevelThree], address, models.MachineLevelFour)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelFour:
		err := machineLevelStepTwo(machineLevel[models.MachineLevelFour], address, models.MachineLevelFive)
		if err != nil {
			mylog.DataLogger.Error().Msgf("machineLevelStepOne machineLevelStepTwo err: %v", err)
		}
	}
}

func machineLevelStepTwo(machineLevel *models.MachineLevel, address *models.Address, countMachineLevelId int64) error {
	// 获取兑换手续费
	sumFee, err := models.SharedRedis().GetAccountConvertSumFee()
	if err != nil {
		return err
	}

	// 获取升级后的达人级别的数量
	count, err := CountAddressByMachineLevelId(countMachineLevelId)
	if err != nil {
		return err
	}

	// 获取实际分红数量
	amount := sumFee.Mul(machineLevel.GlobalFee)
	if count > 0 {
		amount = amount.Div(decimal.NewFromInt(int64(count)))
	}

	// 获取赠送的矿机
	machine, err := GetMachineById(machineLevel.MachineId)
	if err != nil {
		return err
	}

	address.GlobalFee = machineLevel.GlobalFee
	address.MachineLevelId = countMachineLevelId
	return machineLevelStepThree(address, amount, machine)
}

// 更改用户达人级别，进行分红，赠送矿机
func machineLevelStepThree(address *models.Address, amount decimal.Decimal, machine *models.Machine) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	ytlAsset, err := db.GetAccountAssetForUpdate(address.Id, models.AccountCurrencyYtl)
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
		UserId:      address.Id,
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
