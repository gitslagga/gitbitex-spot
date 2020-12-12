package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/utils"
	"github.com/shopspring/decimal"
	"math"
	"strconv"
	"strings"
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
		issueMap[k]["Rate"] = v.Available.Div(total)
		issueMap[k]["Release"] = v.Available.Div(total).Mul(issueReward)
		issueMap[k]["Deduction"] = v.Available.Div(total).Mul(issueReward).Mul(biteRate)
	}

	return issueMap, nil
}

func BackendIssueStart() error {
	issueList, err := BackendIssueList()
	if err != nil {
		return err
	}
	issueConfigs, err := mysql.SharedStore().GetIssueConfigs()
	if err != nil {
		return err
	}

	for _, v := range issueList {
		err = backendIssueStart(int64(v["UserId"].(float64)), v["Deduction"].(decimal.Decimal), issueConfigs)
		if err != nil {
			mylog.Logger.Error().Msgf("BackendIssueStart userId:%v, err:%v", v["UserId"], err)
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

	//TODO Use entire account assets
	accountPools, err := mysql.SharedStore().GetHoldingAccountPool(minHolding)
	if err != nil {
		return nil, err
	}

	var holdingMap = make([]map[string]interface{}, len(accountPools))
	var totalRank decimal.Decimal
	var bestHolding decimal.Decimal
	for k, v := range accountPools {
		holdingMap[k] = utils.StructToMapViaJson(v)

		holdingMap[k]["Rank"] = decimal.NewFromInt(int64(k + 1))
		if k > 0 && v.Available.Equal(accountPools[k-1].Available) {
			holdingMap[k]["Rank"] = holdingMap[k-1]["Rank"]
		}

		totalRank = totalRank.Add(holdingMap[k]["Rank"].(decimal.Decimal))

		if k > 0 && holdingMap[k]["Rank"].(decimal.Decimal).Div(v.Available).
			GreaterThan(holdingMap[k-1]["Rank"].(decimal.Decimal).Div(accountPools[k-1].Available)) {
			bestHolding = v.Available
		}
	}

	for _, val := range holdingMap {
		// 拼团成功，持币排名翻倍
		if models.SharedRedis().ExistsAccountGroupWinTime(int64(val["UserId"].(float64))) {
			val["Rank"] = val["Rank"].(decimal.Decimal).Add(val["Rank"].(decimal.Decimal))
			val["Goal"] = "排名翻倍"
		}

		val["HoldReward"] = holdReward
		val["TotalRank"] = totalRank
		val["Profit"] = val["Rank"].(decimal.Decimal).Div(totalRank).Mul(holdReward).Truncate(8)
	}

	return map[string]interface{}{
		"HoldingMap":  holdingMap,
		"BestHolding": bestHolding,
	}, nil
}

func BackendHoldingStart() error {
	holdingList, err := BackendHoldingList()
	if err != nil {
		return err
	}

	for _, v := range holdingList["HoldingMap"].([]map[string]interface{}) {
		err = backendHoldingStart(int64(v["UserId"].(float64)), v["HoldReward"].(decimal.Decimal), v["Profit"].(decimal.Decimal),
			v["TotalRank"].(decimal.Decimal), v["Rank"].(decimal.Decimal))
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

	accountPool, err := db.GetAccountPoolForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	accountPool.Available = accountPool.Available.Add(number.Mul(decimal.NewFromFloat(1 - models.AccountHoldingShopRate)))
	err = db.UpdateAccountPool(accountPool)
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
		Coin:      models.AccountCurrencyBite,
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

func BackendPromoteList() ([]map[string]interface{}, error) {
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		return nil, err
	}
	promoteReward, err := decimal.NewFromString(configs[models.ConfigPromoteProfit].Value)
	if err != nil {
		return nil, err
	}

	//TODO Use entire account assets
	totalPowerList, err := mysql.SharedStore().GetPromoteAccountPool()
	if err != nil {
		return nil, err
	}

	var parentPower []map[string]interface{}
	var userPower = make(map[int64][]map[string]interface{}, len(totalPowerList))
	var parentIds = make([]string, len(totalPowerList))

	var parentIdInt int64
	for _, v := range totalPowerList {
		parentIds = strings.Split(v.ParentIds, ",")

		for _, parentId := range parentIds {
			parentIdInt, err = strconv.ParseInt(parentId, 10, 64)
			if err != nil {
				mylog.Logger.Error().Msgf("BackendPromoteList ParseInt parentId:%v, err:%v", parentId, err)
				continue
			}

			userPower[parentIdInt] = append(userPower[parentIdInt], map[string]interface{}{
				"UserId":    v.Id,
				"Available": v.Available,
			})
		}
	}

	var maxKey int
	var maxPower decimal.Decimal
	var power decimal.Decimal
	var availableF float64
	var totalPower decimal.Decimal
	for parentId, sonList := range userPower {
		maxKey = 0
		maxPower = sonList[0]["Available"].(decimal.Decimal)
		power = decimal.Zero

		for key, val := range sonList {
			if val["Available"].(decimal.Decimal).GreaterThan(maxPower) {
				maxKey = key
				maxPower = val["Available"].(decimal.Decimal)
			}
		}

		// 算力计算
		for key, val := range sonList {
			if key == maxKey {
				availableF, _ = val["Available"].(decimal.Decimal).Float64()
				userPower[parentId][key]["Power"] = decimal.NewFromFloat(math.Sqrt(math.Sqrt(availableF)))
				power = power.Add(userPower[parentId][key]["Power"].(decimal.Decimal))
				continue
			}

			if val["Available"].(decimal.Decimal).GreaterThan(decimal.NewFromInt(models.AccountPromoteStandard)) {
				userPower[parentId][key]["Power"] = val["Available"].(decimal.Decimal).Add(decimal.NewFromInt(models.AccountPromoteMaxCal))
				power = power.Add(userPower[parentId][key]["Power"].(decimal.Decimal))
			} else if val["Available"].(decimal.Decimal).LessThanOrEqual(decimal.NewFromInt(models.AccountPromoteStandard)) {
				userPower[parentId][key]["Power"] = val["Available"].(decimal.Decimal).Mul(decimal.NewFromInt(models.AccountPromoteMinCal))
				power = power.Add(userPower[parentId][key]["Power"].(decimal.Decimal))
			}
		}

		parentPower = append(parentPower, map[string]interface{}{
			"ParentId": parentId,
			"Power":    power.Truncate(8),
			"Currency": models.AccountCurrencyBite,
			"CountSon": len(sonList),
		})
		totalPower = totalPower.Add(power)
	}

	for _, val := range parentPower {
		val["TotalPower"] = totalPower.Truncate(8)
		val["Profit"] = val["Power"].(decimal.Decimal).Div(totalPower).Mul(promoteReward).Truncate(8)
	}

	return parentPower, nil
}

func BackendPromoteStart() error {
	parentPower, err := BackendPromoteList()
	if err != nil {
		return err
	}
	for _, val := range parentPower {
		err = backendPromoteStart(val["ParentId"].(int64), val["Power"].(decimal.Decimal), val["TotalPower"].(decimal.Decimal),
			val["Profit"].(decimal.Decimal), val["CountSon"].(int))
		if err != nil {
			mylog.Logger.Error().Msgf("BackendPromoteStart userId:%v, err:%v", val["ParentId"], err)
		}
	}

	return nil
}

func backendPromoteStart(userId int64, power, totalPower, profit decimal.Decimal, countSon int) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountPool, err := db.GetAccountPoolForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	accountPool.Available = accountPool.Available.Add(profit.Mul(decimal.NewFromFloat(1 - models.AccountPromoteShopRate)))
	err = db.UpdateAccountPool(accountPool)
	if err != nil {
		return err
	}

	shopAsset, err := db.GetAccountShopForUpdate(userId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	shopAsset.Available = shopAsset.Available.Add(profit.Mul(decimal.NewFromFloat(models.AccountPromoteShopRate)))
	err = db.UpdateAccountShop(shopAsset)
	if err != nil {
		return err
	}

	err = db.AddAddressPromote(&models.AddressPromote{
		UserId:     userId,
		Coin:       models.AccountCurrencyBite,
		Number:     profit,
		Power:      power,
		TotalPower: totalPower,
		CountSon:   countSon,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}
