package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/utils"
	"github.com/shopspring/decimal"
)

func BackendIssueList() ([]map[string]interface{}, error) {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return nil, err
	}

	issueReward, err := decimal.NewFromString(configs[models.ConfigIssueReward].Value)
	if err != nil {
		return nil, err
	}
	biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
	if err != nil {
		return nil, err
	}

	accountAssets, err := mysql.SharedStore().GetIssueAccountAsset()
	if err != nil {
		return nil, err
	}
	total, err := mysql.SharedStore().SumIssueAccountAsset()
	if err != nil {
		return nil, err
	}

	var issueMap = make([]map[string]interface{}, len(accountAssets))
	for k, v := range accountAssets {
		issueMap[k] = utils.StructToMapViaJson(v)
		issueMap[k]["rate"] = v.Available.Div(total)
		issueMap[k]["release"] = v.Available.Div(total).Mul(issueReward)
		issueMap[k]["deduction"] = v.Available.Div(total).Mul(issueReward).Mul(biteRate)
	}

	return issueMap, nil
}

func BackendIssueStart() error {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	issueReward, err := decimal.NewFromString(configs[models.ConfigIssueReward].Value)
	if err != nil {
		return err
	}
	biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
	if err != nil {
		return err
	}

	accountAssets, err := mysql.SharedStore().GetIssueAccountAsset()
	if err != nil {
		return err
	}
	total, err := mysql.SharedStore().SumIssueAccountAsset()
	if err != nil {
		return err
	}
	issueConfigs, err := mysql.SharedStore().GetIssueConfigs()
	if err != nil {
		return err
	}

	for _, v := range accountAssets {
		err = backendIssueStart(v.UserId, v.Available.Div(total).Mul(issueReward).Mul(biteRate), issueConfigs)
		if err != nil {
			mylog.Logger.Error().Msgf("BackendIssueStart userId:%v, err:%v", v.UserId, err)
		}
	}

	return nil
}

func backendIssueStart(userId int64, deduction decimal.Decimal, config []*models.IssueConfig) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	coinAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyUsdt)
	if err != nil {
		return err
	}

	if coinAsset.Available.LessThan(deduction) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	coinAsset.Available = coinAsset.Available.Sub(deduction)
	err = db.UpdateAccountAsset(coinAsset)
	if err != nil {
		return err
	}

	err = db.AddIssue(&models.Issue{
		UserId:       userId,
		Coin:         models.AccountCurrencyUsdt,
		Number:       deduction,
		Remain:       deduction,
		Count:        0,
		ReleaseOne:   config[models.IssueFirstEveryMonth].Number,
		ReleaseTwo:   config[models.IssueSecondEveryMonth].Number,
		ReleaseThree: config[models.IssueThreeEveryMonth].Number,
		ReleaseFour:  config[models.IssueFourEveryMonth].Number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func GetLastAddressHolding() (*models.AddressHolding, error) {
	return mysql.SharedStore().GetLastAddressHolding()
}

func GetLastAddressPromote() (*models.AddressPromote, error) {
	return mysql.SharedStore().GetLastAddressPromote()
}

func BackendHoldingList() (map[string]interface{}, error) {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return nil, err
	}

	holdReward, err := decimal.NewFromString(configs[models.ConfigHoldCoinProfit].Value)
	if err != nil {
		return nil, err
	}
	minHolding, err := decimal.NewFromString(configs[models.ConfigMinHolding].Value)
	if err != nil {
		return nil, err
	}

	accountAssets, err := mysql.SharedStore().GetHoldingAccountAsset(minHolding)
	if err != nil {
		return nil, err
	}

	var holdingMap = make([]map[string]interface{}, len(accountAssets))
	var totalRank decimal.Decimal
	var bestHolding decimal.Decimal
	for k, v := range accountAssets {
		holdingMap[k] = utils.StructToMapViaJson(v)

		holdingMap[k]["rank"] = decimal.NewFromInt(int64(k + 1))
		if k > 0 && v.Available.Equal(accountAssets[k-1].Available) {
			holdingMap[k]["rank"] = holdingMap[k-1]["rank"]
		}

		totalRank = totalRank.Add(holdingMap[k]["rank"].(decimal.Decimal))

		if k > 0 && holdingMap[k]["rank"].(decimal.Decimal).Div(v.Available).
			GreaterThan(holdingMap[k-1]["rank"].(decimal.Decimal).Div(accountAssets[k-1].Available)) {
			bestHolding = v.Available
		}
	}

	for k, _ := range holdingMap {
		holdingMap[k]["profit"] = holdingMap[k]["rank"].(decimal.Decimal).Div(totalRank).Mul(holdReward)
	}

	return map[string]interface{}{
		"holding_map":  holdingMap,
		"best_holding": bestHolding,
	}, nil
}

func BackendHoldingStart() error {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	holdReward, err := decimal.NewFromString(configs[models.ConfigHoldCoinProfit].Value)
	if err != nil {
		return err
	}
	minHolding, err := decimal.NewFromString(configs[models.ConfigMinHolding].Value)
	if err != nil {
		return err
	}

	accountAssets, err := mysql.SharedStore().GetHoldingAccountAsset(minHolding)
	if err != nil {
		return err
	}

	var holdingMap = make([]map[string]interface{}, len(accountAssets))
	var totalRank decimal.Decimal
	for k, v := range accountAssets {
		holdingMap[k] = utils.StructToMapViaJson(v)

		holdingMap[k]["rank"] = decimal.NewFromInt(int64(k + 1))
		if k > 0 && v.Available.Equal(accountAssets[k-1].Available) {
			holdingMap[k]["rank"] = holdingMap[k-1]["rank"]
		}
		totalRank = totalRank.Add(holdingMap[k]["rank"].(decimal.Decimal))
	}

	for k, v := range holdingMap {
		holdingMap[k]["profit"] = holdingMap[k]["rank"].(decimal.Decimal).Div(totalRank).Mul(holdReward)
		err = backendHoldingStart(int64(v["UserId"].(float64)), holdReward, holdingMap[k]["profit"].(decimal.Decimal),
			totalRank, holdingMap[k]["rank"].(decimal.Decimal))
		if err != nil {
			mylog.Logger.Error().Msgf("backendHoldingStart userId:%v, err:%v", v["userId"], err)
		}
	}

	return nil
}

func backendHoldingStart(userId int64, totalNum, number, totalRank, rank decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	accountAsset.Available = accountAsset.Available.Add(number.Mul(decimal.NewFromFloat(1 - models.AccountHoldingShopRate)))
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	shopAsset, err := db.GetAccountShopForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	shopAsset.Available = shopAsset.Available.Add(number.Mul(decimal.NewFromFloat(models.AccountHoldingShopRate)))
	err = db.UpdateAccountShop(shopAsset)
	if err != nil {
		return err
	}

	err = db.AddAddressHolding(&models.AddressHolding{
		UserId:    userId,
		Coin:      "BITE",
		TotalNum:  totalNum,
		Number:    number,
		TotalRank: int(totalRank.IntPart()),
		Rank:      int(rank.IntPart()),
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func BackendPromoteList() (map[string]interface{}, error) {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return nil, err
	}

	holdReward, err := decimal.NewFromString(configs[models.ConfigHoldCoinProfit].Value)
	if err != nil {
		return nil, err
	}
	minHolding, err := decimal.NewFromString(configs[models.ConfigMinHolding].Value)
	if err != nil {
		return nil, err
	}

	accountAssets, err := mysql.SharedStore().GetHoldingAccountAsset(minHolding)
	if err != nil {
		return nil, err
	}

	var holdingMap = make([]map[string]interface{}, len(accountAssets))
	var totalRank decimal.Decimal
	var bestHolding decimal.Decimal
	for k, v := range accountAssets {
		holdingMap[k] = utils.StructToMapViaJson(v)

		holdingMap[k]["rank"] = decimal.NewFromInt(int64(k + 1))
		if k > 0 && v.Available.Equal(accountAssets[k-1].Available) {
			holdingMap[k]["rank"] = holdingMap[k-1]["rank"]
		}

		totalRank = totalRank.Add(holdingMap[k]["rank"].(decimal.Decimal))

		if k > 0 && holdingMap[k]["rank"].(decimal.Decimal).Div(v.Available).
			GreaterThan(holdingMap[k-1]["rank"].(decimal.Decimal).Div(accountAssets[k-1].Available)) {
			bestHolding = v.Available
		}
	}

	for k, _ := range holdingMap {
		holdingMap[k]["profit"] = holdingMap[k]["rank"].(decimal.Decimal).Div(totalRank).Mul(holdReward)
	}

	return map[string]interface{}{
		"holding_map":  holdingMap,
		"best_holding": bestHolding,
	}, nil
}

func BackendPromoteStart() error {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return err
	}

	holdReward, err := decimal.NewFromString(configs[models.ConfigHoldCoinProfit].Value)
	if err != nil {
		return err
	}
	minHolding, err := decimal.NewFromString(configs[models.ConfigMinHolding].Value)
	if err != nil {
		return err
	}

	accountAssets, err := mysql.SharedStore().GetHoldingAccountAsset(minHolding)
	if err != nil {
		return err
	}

	var holdingMap = make([]map[string]interface{}, len(accountAssets))
	var totalRank decimal.Decimal
	for k, v := range accountAssets {
		holdingMap[k] = utils.StructToMapViaJson(v)

		holdingMap[k]["rank"] = decimal.NewFromInt(int64(k + 1))
		if k > 0 && v.Available.Equal(accountAssets[k-1].Available) {
			holdingMap[k]["rank"] = holdingMap[k-1]["rank"]
		}
		totalRank = totalRank.Add(holdingMap[k]["rank"].(decimal.Decimal))
	}

	for k, v := range holdingMap {
		holdingMap[k]["profit"] = holdingMap[k]["rank"].(decimal.Decimal).Div(totalRank).Mul(holdReward)
		err = backendPromoteStart(int64(v["UserId"].(float64)), holdingMap[k]["profit"].(decimal.Decimal))
		if err != nil {
			mylog.Logger.Error().Msgf("backendHoldingStart userId:%v, err:%v", v["userId"], err)
		}
	}

	return nil
}

func backendPromoteStart(userId int64, profit decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	accountAsset.Available = accountAsset.Available.Add(profit.Mul(decimal.NewFromFloat(1 - models.AccountHoldingShopRate)))
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	shopAsset, err := db.GetAccountShopForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	shopAsset.Available = shopAsset.Available.Add(profit.Mul(decimal.NewFromFloat(models.AccountHoldingShopRate)))
	err = db.UpdateAccountShop(shopAsset)
	if err != nil {
		return err
	}

	err = db.AddAddressHolding(&models.AddressHolding{
		UserId: userId,
		Coin:   "BITE",
		Number: profit,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
