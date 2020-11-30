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
