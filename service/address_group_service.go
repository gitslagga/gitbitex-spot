package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/utils"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetAddressGroupByUserId(userId int64) ([]*models.AddressGroup, error) {
	return mysql.SharedStore().GetAddressGroupByUserId(userId)
}

func AddAddressGroup(group *models.AddressGroup) error {
	return mysql.SharedStore().AddAddressGroup(group)
}

func UpdateAddressGroup(group *models.AddressGroup) error {
	return mysql.SharedStore().UpdateAddressGroup(group)
}

func AddressGroup(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:

		// TODO 兑换手续费，扫描手续费，节点收益 定时任务
		//sumFee, err := mysql.SharedStore().GetMachineConvertSumFee()
		//if err == nil {
		//	_ = models.SharedRedis().SetAccountGroupSumNum(sumFee)
		//}

		return addressGroup(address, coin, decimal.NewFromInt(models.AccountGroupPlayUsdt), decimal.NewFromInt(models.AccountGroupRefundUsdt))
	case models.AccountGroupCurrencyBite:
		return addressGroup(address, coin, decimal.NewFromInt(models.AccountGroupPlayBite), decimal.NewFromInt(models.AccountGroupRefundBite))
	}

	return nil
}

func addressGroup(address *models.Address, coin string, number, refund decimal.Decimal) error {
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

	addressGroup, err := db.GetAddressGroupForUpdate(coin)
	if err != nil {
		return err
	}
	if addressGroup == nil || addressGroup.Count >= models.AccountGroupPersonNumber {
		// 创建一个新的拼团
		err = db.AddAddressGroup(&models.AddressGroup{
			UserId:  address.Id,
			Coin:    coin,
			OrderSN: utils.GetOrderSN(),
			Number:  number,
			Status:  models.AccountGroupDefault,
			Count:   1,
		})
	} else {
		addressGroupExists, err := db.GetAddressGroupByUserIdOrderSN(address.Id, addressGroup.OrderSN)
		if err != nil {
			return err
		}
		if addressGroupExists != nil {
			return errors.New("本轮已经加入|The current round has been joined")
		}

		// 使用旧的拼团
		addressGroup.Count++
		err = db.AddAddressGroup(&models.AddressGroup{
			UserId:  address.Id,
			Coin:    addressGroup.Coin,
			OrderSN: addressGroup.OrderSN,
			Number:  number,
			Status:  models.AccountGroupDefault,
			Count:   addressGroup.Count,
		})

		if addressGroup.Count >= models.AccountGroupPersonNumber {
			// 拼团成功
			addressGroups, err := db.GetAddressGroupsByOrderSN(addressGroup.OrderSN)
			if err != nil {
				return err
			}

			goal := rand.Intn(models.AccountGroupPersonNumber)
			addressGroups[goal].Status = models.AccountGroupWinning
			err = db.UpdateAddressGroup(addressGroups[goal])
			if err != nil {
				return err
			}

			// 持币排名翻倍，再次成功持续时间翻倍
			ttl := models.SharedRedis().TtlAccountGroupWinTime(addressGroups[goal].UserId)

			var exp time.Duration
			if ttl.Seconds() == -2 {
				exp = models.AccountGroupIncreaseDay * 24 * time.Hour
			} else {
				exp = ttl + models.AccountGroupIncreaseDay*24*time.Hour
			}

			err = models.SharedRedis().SetAccountGroupWinTime(addressGroups[goal].UserId, exp)
			if err != nil {
				mylog.Logger.Error().Msgf("addressGroup SetAccountGroupWinTime err:%v", err)
			}

			// 拼团失败，返还数量
			for key, val := range addressGroups {
				if key == goal {
					continue
				}

				accountAsset, err = db.GetAccountAssetForUpdate(val.UserId, val.Coin)
				if err != nil {
					return err
				}
				accountAsset.Available = accountAsset.Available.Add(refund)
				err = db.UpdateAccountAsset(accountAsset)
				if err != nil {
					return err
				}
			}
		}
	}

	// 直推收益
	if address.ParentId != 0 {
		accountAsset, err = db.GetAccountAssetForUpdate(address.ParentId, coin)
		if err != nil {
			return err
		}
		accountAsset.Available = accountAsset.Available.Add(number.Mul(decimal.NewFromFloat(models.AccountGroupDirectRate)))
		err = db.UpdateAccountAsset(accountAsset)
		if err != nil {
			return err
		}
	}

	// 间推收益
	if address.ParentIds != "" {
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
	}

	return db.CommitTx()
}

func AddressDelegate(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:
		if address.GroupUsdt == models.AccountGroupDelegateMode {
			return errors.New("已经质押过|You have been delegated")
		}

		address.GroupUsdt = models.AccountGroupDelegateMode
		return addressDelegate(address, coin, decimal.NewFromInt(models.AccountGroupStakeUsdt))
	case models.AccountGroupCurrencyBite:
		if address.GroupBite == models.AccountGroupDelegateMode {
			return errors.New("已经质押过|You have been delegated")
		}

		address.GroupBite = models.AccountGroupDelegateMode
		return addressDelegate(address, coin, decimal.NewFromInt(models.AccountGroupStakeBite))
	}

	return nil
}

func addressDelegate(address *models.Address, coin string, number decimal.Decimal) error {
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

func AddressRelease(address *models.Address, coin string) error {
	switch coin {
	case models.AccountGroupCurrencyUsdt:
		if address.GroupUsdt == models.AccountGroupReleaseMode {
			return errors.New("已经释放过|You have been released")
		}

		address.GroupUsdt = models.AccountGroupReleaseMode
		return addressRelease(address, coin, decimal.NewFromInt(models.AccountGroupStakeUsdt))
	case models.AccountGroupCurrencyBite:
		if address.GroupBite == models.AccountGroupReleaseMode {
			return errors.New("已经释放过|You have been released")
		}

		address.GroupBite = models.AccountGroupReleaseMode
		return addressRelease(address, coin, decimal.NewFromInt(models.AccountGroupStakeBite))
	}

	return nil
}

func addressRelease(address *models.Address, coin string, number decimal.Decimal) error {
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
