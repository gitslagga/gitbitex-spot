package tasks

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"time"
)

// 认购释放任务
func StartIssueRelease() {
	IssueRelease()

	t := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-t.C:
			IssueRelease()
		}
	}
}

func IssueRelease() {
	//获取BITE兑换USDT费率
	configs, err := mysql.SharedStore().GetConfigs()
	if err != nil {
		mylog.DataLogger.Error().Msgf("IssueRelease GetConfigs err: %v", err)
		return
	}
	biteRate, err := decimal.NewFromString(configs[models.ConfigBiteConvertUsdt].Value)
	if err != nil {
		mylog.DataLogger.Error().Msgf("IssueRelease biteRate err: %v", err)
		return
	}
	if biteRate.LessThanOrEqual(decimal.Zero) {
		mylog.DataLogger.Error().Msgf("IssueRelease BITE convert USDT price error")
		return
	}

	issueUsedList, err := mysql.SharedStore().GetIssueUsedList()
	if err == nil {
		for _, v := range issueUsedList {
			//获取认购最后一条释放记录
			issueLog, err := mysql.SharedStore().GetLastIssueLog(v.Id)
			if err != nil {
				mylog.DataLogger.Error().Msgf("IssueRelease GetLastIssueLog err: %v", err)
				continue
			}
			if issueLog != nil {
				currentTime := time.Now()
				startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 00, 00, 00, 00, currentTime.Location())
				endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location())

				if issueLog.CreatedAt.After(startTime) && issueLog.CreatedAt.Before(endTime) {
					continue
				}
			}

			//获取实际应得的BITE数量
			var number decimal.Decimal
			if time.Now().Before(v.CreatedAt.Add(90 * 24 * time.Hour)) {
				number = v.Number.Mul(v.ReleaseOne)
			} else if time.Now().Before(v.CreatedAt.Add(180 * 24 * time.Hour)) {
				number = v.Number.Mul(v.ReleaseTwo)
			} else if time.Now().Before(v.CreatedAt.Add(270 * 24 * time.Hour)) {
				number = v.Number.Mul(v.ReleaseThree)
			} else {
				number = v.Number.Mul(v.ReleaseFour)
			}
			if v.Remain.LessThan(number) {
				number = v.Remain
			}

			err = issueRelease(v, number, number.Div(biteRate))
			if err != nil {
				mylog.DataLogger.Error().Msgf("IssueRelease machineRelease err: %v", err)
			}
		}
	}
}

func issueRelease(issue *models.Issue, number, biteNumber decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	issue.Count++
	issue.Remain = issue.Remain.Sub(number)
	err = db.UpdateIssue(issue)
	if err != nil {
		return err
	}

	err = db.AddIssueLog(&models.IssueLog{
		UserId:  issue.UserId,
		IssueId: issue.Id,
		Number:  biteNumber,
	})
	if err != nil {
		return err
	}

	addressAsset, err := db.GetAccountAssetForUpdate(issue.UserId, models.AccountCurrencyBite)
	if err != nil {
		return err
	}

	addressAsset.Available = addressAsset.Available.Add(biteNumber)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
