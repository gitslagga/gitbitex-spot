package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetGroupById(groupId int64) (*models.Group, error) {
	return mysql.SharedStore().GetGroupById(groupId)
}

func GetGroupByUserIdCoin(userId int64, coin string) (*models.Group, error) {
	return mysql.SharedStore().GetGroupByUserIdCoin(userId, coin)
}

func GetGroupByCoin(coin string, before, after, limit int64) ([]*models.Group, error) {
	return mysql.SharedStore().GetGroupByCoin(coin, before, after, limit)
}

func GetGroupLogByUserId(userId, before, after, limit int64) ([]*models.GroupLog, error) {
	return mysql.SharedStore().GetGroupLogByUserId(userId, before, after, limit)
}

func GetGroupLogPublicity(before, after, limit int64) ([]*models.GroupLog, error) {
	return mysql.SharedStore().GetGroupLogPublicity(before, after, limit)
}

func GetGroupLogByGroupIdUserId(groupId, userId int64) (*models.GroupLog, error) {
	return mysql.SharedStore().GetGroupLogByGroupIdUserId(groupId, userId)
}

func GetAddressReleaseByUserId(userId, before, after, limit int64) ([]*models.AddressRelease, error) {
	return mysql.SharedStore().GetAddressReleaseByUserId(userId, before, after, limit)
}

func GroupPublish(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:
		return groupPublish(address, coin, decimal.NewFromInt(models.AccountGroupPlayUsdt), decimal.NewFromInt(models.AccountGroupRefundUsdt))
	case models.AccountGroupCurrencyBite:
		return groupPublish(address, coin, decimal.NewFromInt(models.AccountGroupPlayBite), decimal.NewFromInt(models.AccountGroupRefundBite))
	}

	return nil
}

func groupPublish(address *models.Address, coin string, number, refund decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(address.Id, coin)
	if err != nil {
		return err
	}
	if accountAsset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	accountAsset.Available = accountAsset.Available.Sub(number)
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	// 创建一个新的拼团
	group := &models.Group{
		UserId:    address.Id,
		Address:   address.Address,
		Coin:      coin,
		Fee:       decimal.Zero,
		Number:    number,
		Refund:    refund,
		JoinCount: 1,
		Count:     models.AccountGroupPersonNumber,
		Status:    models.AccountGroupProcess,
	}
	err = db.AddGroup(group)
	if err != nil {
		return err
	}

	err = db.AddGroupLog(&models.GroupLog{
		GroupId: group.Id,
		UserId:  address.Id,
		Address: address.Address,
		Coin:    coin,
		Number:  number,
		Status:  models.AccountGroupLogProcess,
	})
	if err != nil {
		return err
	}

	if address.ParentId != 0 && address.ParentIds != "" {
		err = parentProfit(address, db, number, coin)
		if err != nil {
			return err
		}
	}

	return db.CommitTx()
}

// 直推收益，间推收益
func parentProfit(address *models.Address, db models.Store, number decimal.Decimal, coin string) error {
	accountAsset, err := db.GetAccountAssetForUpdate(address.ParentId, coin)
	if err != nil {
		return err
	}
	accountAsset.Available = accountAsset.Available.Add(number.Mul(decimal.NewFromFloat(models.AccountGroupDirectRate)))
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	parentIds := strings.Split(address.ParentIds, ",")
	for i := 0; i < len(parentIds)-1; i++ {
		parentId, err := strconv.ParseInt(parentIds[i], 10, 64)
		if err != nil {
			return err
		}

		accountAsset, err = db.GetAccountAssetForUpdate(parentId, coin)
		if err != nil {
			return err
		}
		accountAsset.Available = accountAsset.Available.Add(number.Mul(decimal.NewFromFloat(models.AccountGroupIndirectRate)))
		err = db.UpdateAccountAsset(accountAsset)
		if err != nil {
			return err
		}
	}
	return nil
}

func GroupJoin(address *models.Address, group *models.Group) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	accountAsset, err := db.GetAccountAssetForUpdate(address.Id, group.Coin)
	if err != nil {
		return err
	}
	if accountAsset.Available.LessThan(group.Number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	accountAsset.Available = accountAsset.Available.Sub(group.Number)
	err = db.UpdateAccountAsset(accountAsset)
	if err != nil {
		return err
	}

	// 加入拼团
	group.JoinCount++
	if group.JoinCount >= group.Count {
		group.Status = models.AccountGroupFinish
	}
	err = db.UpdateGroup(group)
	if err != nil {
		return err
	}

	err = db.AddGroupLog(&models.GroupLog{
		GroupId: group.Id,
		UserId:  address.Id,
		Address: address.Address,
		Coin:    group.Coin,
		Number:  group.Number,
		Status:  models.AccountGroupLogProcess,
	})
	if err != nil {
		return err
	}

	if group.JoinCount >= group.Count {
		rand.Seed(time.Now().UnixNano())
		goal := rand.Intn(group.Count)

		// 持币排名翻倍，再次成功持续时间翻倍
		groupLogs, err := db.GetGroupLogByGroupId(group.Id)
		if err != nil {
			return err
		}
		ttl := models.SharedRedis().TtlAccountGroupWinTime(groupLogs[goal].UserId)

		var exp time.Duration
		if ttl.Seconds() == -2 {
			exp = models.AccountGroupIncreaseDay * 24 * time.Hour
		} else {
			exp = ttl + models.AccountGroupIncreaseDay*24*time.Hour
		}

		err = models.SharedRedis().SetAccountGroupWinTime(groupLogs[goal].UserId, exp)
		if err != nil {
			mylog.Logger.Error().Msgf("GroupJoin SetAccountGroupWinTime err:%v", err)
		}

		// 拼团成功
		groupLogs[goal].Status = models.AccountGroupLogSuccess
		err = db.UpdateGroupLog(groupLogs[goal])
		if err != nil {
			return err
		}

		// 拼团失败，返还数量
		for key, val := range groupLogs {
			if key == goal {
				continue
			}

			val.Status = models.AccountGroupLogFailed
			err = db.UpdateGroupLog(val)
			if err != nil {
				return err
			}

			accountAsset, err = db.GetAccountAssetForUpdate(val.UserId, val.Coin)
			if err != nil {
				return err
			}
			accountAsset.Available = accountAsset.Available.Add(group.Refund)
			err = db.UpdateAccountAsset(accountAsset)
			if err != nil {
				return err
			}
		}
	}

	if address.ParentId != 0 && address.ParentIds != "" {
		err = parentProfit(address, db, group.Number, group.Coin)
		if err != nil {
			return err
		}
	}

	return db.CommitTx()
}

func GroupDelegate(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:
		if address.GroupUsdt == models.AccountGroupDelegateMode {
			return errors.New("已经质押过|You have been delegated")
		}

		address.GroupUsdt = models.AccountGroupDelegateMode
		return groupDelegate(address, coin, decimal.NewFromInt(models.AccountGroupStakeUsdt))
	case models.AccountGroupCurrencyBite:
		if address.GroupBite == models.AccountGroupDelegateMode {
			return errors.New("已经质押过|You have been delegated")
		}

		address.GroupBite = models.AccountGroupDelegateMode
		return groupDelegate(address, coin, decimal.NewFromInt(models.AccountGroupStakeBite))
	}

	return nil
}

func groupDelegate(address *models.Address, coin string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(address.Id, coin)
	if err != nil {
		return err
	}

	if addressAsset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}

	addressAsset.Available = addressAsset.Available.Sub(number)
	addressAsset.Hold = addressAsset.Hold.Add(number)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func GroupRelease(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:
		if address.GroupUsdt == models.AccountGroupReleaseMode {
			return errors.New("已经释放过|You have been released")
		}

		address.GroupUsdt = models.AccountGroupReleaseMode
		return groupRelease(address, coin, decimal.NewFromInt(models.AccountGroupStakeUsdt))
	case models.AccountGroupCurrencyBite:
		if address.GroupBite == models.AccountGroupReleaseMode {
			return errors.New("已经释放过|You have been released")
		}

		address.GroupBite = models.AccountGroupReleaseMode
		return groupRelease(address, coin, decimal.NewFromInt(models.AccountGroupStakeBite))
	}

	return nil
}

func groupRelease(address *models.Address, coin string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	addressAsset, err := db.GetAccountAssetForUpdate(address.Id, coin)
	if err != nil {
		return err
	}

	addressAsset.Available = addressAsset.Available.Add(number)
	addressAsset.Hold = addressAsset.Hold.Sub(number)
	err = db.UpdateAccountAsset(addressAsset)
	if err != nil {
		return err
	}

	err = db.UpdateAddress(address)
	if err != nil {
		return err
	}

	return db.CommitTx()
}
