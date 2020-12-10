package models

import (
	"fmt"
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

const (
	AccountConvertSumFee = "account_convert_sum_fee"
	AccountScanSumFee    = "account_scan_sum_fee"
	AccountGroupSumNum   = "account_group_sum_num"

	EthLatestHeightEth  = "wallet_latest_height_eth"
	AccountGroupWinTime = "account_group_win_time"
)

var redisClient *redis.Client
var redisOnce sync.Once

type box struct {
	redis *redis.Client
}

func SharedRedis() *box {
	redisOnce.Do(func() {
		gbeConfig := conf.GetConfig()

		redisClient = redis.NewClient(&redis.Options{
			Addr:     gbeConfig.Redis.Addr,
			Password: gbeConfig.Redis.Password,
			DB:       0,
		})
	})
	return &box{redis: redisClient}
}

func (b *box) SetAccountConvertSumFee(sumFee decimal.Decimal, exp time.Duration) error {
	sumFeeF, _ := sumFee.Float64()
	err := b.redis.Set(AccountConvertSumFee, sumFeeF, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) ExistsAccountConvertSumFee() bool {
	exists := b.redis.Exists(AccountConvertSumFee).Val()
	if exists == 0 {
		return false
	}

	return true
}

func (b *box) SetAccountScanSumFee(sumFee decimal.Decimal, exp time.Duration) error {
	sumFeeF, _ := sumFee.Float64()
	err := b.redis.Set(AccountScanSumFee, sumFeeF, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) ExistsAccountScanSumFee() bool {
	exists := b.redis.Exists(AccountScanSumFee).Val()
	if exists == 0 {
		return false
	}

	return true
}

func (b *box) SetEthLatestHeight(height uint64) error {
	err := b.redis.Set(EthLatestHeightEth, height, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) GetEthLatestHeight() (uint64, error) {
	height, err := b.redis.Get(EthLatestHeightEth).Uint64()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return height, nil
}

func (b *box) SetAccountGroupSumNum(sumNum decimal.Decimal, exp time.Duration) error {
	sumFeeF, _ := sumNum.Float64()
	err := b.redis.Set(AccountGroupSumNum, sumFeeF, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) ExistsAccountGroupSumNum() bool {
	exists := b.redis.Exists(AccountGroupSumNum).Val()
	if exists == 0 {
		return false
	}

	return true
}

func (b *box) SetAccountGroupWinTime(userId int64, exp time.Duration) error {
	err := b.redis.Set(fmt.Sprintf("%s_%v", AccountGroupWinTime, userId), userId, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) TtlAccountGroupWinTime(userId int64) time.Duration {
	return b.redis.TTL(fmt.Sprintf("%s_%v", AccountGroupWinTime, userId)).Val()
}

func (b *box) ExistsAccountGroupWinTime(userId int64) bool {
	exists := b.redis.Exists(fmt.Sprintf("%s_%v", AccountGroupWinTime, userId)).Val()
	if exists == 0 {
		return false
	}

	return true
}
