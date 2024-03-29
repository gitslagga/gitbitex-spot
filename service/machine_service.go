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
	count, err := mysql.SharedStore().CountMachineAddressUsed(address.Id, machine.Id)
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

	//获取YTL兑换USDT费率
	ytlRate, err := decimal.NewFromString(configs[models.ConfigYtlConvertUsdt].Value)
	if err != nil {
		return err
	}
	if ytlRate.LessThanOrEqual(decimal.Zero) {
		return errors.New("YTL兑换USDT价格错误|YTL convert USDT price error")
	}

	var amount decimal.Decimal
	switch currency {
	case models.AccountCurrencyYtl:
		amount = machine.Number.Div(ytlRate)
	case models.AccountCurrencyBite:
		biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
		if err != nil {
			return err
		}
		if biteRate.LessThanOrEqual(decimal.Zero) {
			return errors.New("BITE兑换USDT价格错误|BITE convert USDT price error")
		}
		amount = machine.Number.Div(biteRate)
	case models.AccountCurrencyUsdt:
		amount = machine.Number
	default:
		return errors.New("无效的币种|Invalid of currency")
	}

	return buyMachine(address, machine, currency, amount, ytlRate)
}

func buyMachine(address *models.Address, machine *models.Machine, currency string, amount, ytlRate decimal.Decimal) error {
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

	if address.ParentId != 0 {
		//增加上级直推奖励
		parentAddressAsset, err := db.GetAccountAssetForUpdate(address.ParentId, models.AccountCurrencyYtl)
		if err != nil {
			return err
		}

		parentAddressAsset.Available = parentAddressAsset.Available.Add(machine.Number.Mul(machine.Invite).Div(ytlRate))
		err = db.UpdateAccountAsset(parentAddressAsset)
		if err != nil {
			return err
		}

	}
	if address.ParentId != 0 && machine.Id >= models.MachineEffectInvite {
		//增加上级有效账户，更改上级糖果兑换手续费
		parentAddress, err := GetAddressById(address.ParentId)
		if err != nil {
			return err
		}
		parentAddress.InviteNum++

		configs, err := db.GetMachineConfigs()
		if err != nil {
			return err
		}
		for _, v := range configs {
			if parentAddress.InviteNum >= v.InviteNum {
				parentAddress.ConvertFee = v.ConvertFee
			}
		}

		err = db.UpdateAddress(parentAddress)
		if err != nil {
			return err
		}
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

	var activeNum int
	var totalNum int
	for i := 0; i < len(sonAddress); i++ {
		activeNum, err = getTotalActiveNumber(sonAddress[i].Id)
		if err != nil {
			return 0, err
		}

		totalNum += sonAddress[i].ActiveNum
		totalNum += activeNum
	}

	return totalNum, nil
}

func StartMachineLevel(address *models.Address) {
	machineLevel, err := GetMachineLevel()
	if err != nil {
		mylog.Logger.Error().Msgf("StartMachineLevel GetMachineLevel err: %v", err)
		return
	}

	switch address.MachineLevelId {
	case models.MachineLevelZero:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelZero])
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelZero], address, models.MachineLevelOne)
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelOne:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelOne])
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelOne], address, models.MachineLevelTwo)
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelTwo:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelTwo])
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelTwo], address, models.MachineLevelThree)
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelThree:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelThree])
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelThree], address, models.MachineLevelFour)
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel machineLevelStepTwo err: %v", err)
		}
	case models.MachineLevelFour:
		valid, err := getConditionMachineLevel(address, machineLevel[models.MachineLevelFour])
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel getConditionMachineLevel err: %v", err)
			return
		}
		if !valid {
			return
		}

		err = machineLevelStepTwo(machineLevel[models.MachineLevelFour], address, models.MachineLevelFive)
		if err != nil {
			mylog.Logger.Error().Msgf("StartMachineLevel machineLevelStepTwo err: %v", err)
		}
	}
}

func machineLevelStepTwo(machineLevel *models.MachineLevel, address *models.Address, machineLevelId int64) error {
	// 获取赠送的矿机
	machine, err := GetMachineById(machineLevel.MachineId)
	if err != nil {
		return err
	}

	address.GlobalFee = machineLevel.GlobalFee
	address.MachineLevelId = machineLevelId
	return machineLevelStepThree(address, machine)
}

// 更改用户达人级别，赠送矿机
func machineLevelStepThree(address *models.Address, machine *models.Machine) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	err = db.UpdateAddress(address)
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
