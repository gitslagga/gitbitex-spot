package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"time"
)

// YTL兑换BITE资金池任务
func StartConvertRelease() {
	ConvertRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			ConvertRelease()
		}
	}
}

func ConvertRelease() {
	// 查看当天是否释放过
	if models.SharedRedis().ExistsAccountConvertSumFee() {
		return
	}

	// 获取兑换手续费
	sumFee, err := mysql.SharedStore().GetMachineConvertSumFee()
	if err != nil {
		mylog.DataLogger.Error().Msgf("ConvertRelease GetMachineConvertSumFee err: %v", err)
		return
	}

	machineLevel, err := mysql.SharedStore().GetMachineLevel()
	if err != nil {
		mylog.Logger.Error().Msgf("ConvertRelease GetMachineLevel err: %v", err)
		return
	}

	for _, val := range machineLevel {
		// 获取升级后的达人级别的数量
		count, err := mysql.SharedStore().CountAddressByMachineLevelId(val.Id)
		if err != nil {
			return
		}
		if count <= 0 {
			continue
		}

		addresses, err := mysql.SharedStore().GetAddressByMachineLevelId(val.Id)
		if err != nil {
			return
		}
		if addresses == nil {
			continue
		}

		// 获取实际分红数量
		for _, address := range addresses {
			err = convertRelease(address, sumFee.Mul(val.GlobalFee).Div(decimal.NewFromInt(int64(count))))
			if err != nil {
				mylog.Logger.Error().Msgf("ConvertRelease convertRelease err: %v", err)
			}
		}
	}

	// 应该给予资金池足够的释放时间
	_ = models.SharedRedis().SetAccountConvertSumFee(sumFee, time.Hour*23)
}

func convertRelease(address *models.Address, amount decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	ytlAsset, err := db.GetAccountAssetForUpdate(address.Id, models.AccountCurrencyYtl)
	if err != nil {
		return err
	}

	ytlAsset.Available = ytlAsset.Available.Add(amount)
	err = db.UpdateAccountAsset(ytlAsset)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
