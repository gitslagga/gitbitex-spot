package models

import (
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

const (
	AccountConvertSumFee = "account_convert_sum_fee"
	AccountScanSumFee    = "account_scan_sum_fee"

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

func (b *box) SetMachineConvertSumFee(sumFee decimal.Decimal, exp time.Duration) error {
	sumFeeF, _ := sumFee.Float64()
	err := b.redis.Set(AccountConvertSumFee, sumFeeF, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) GetAccountConvertSumFee() (decimal.Decimal, error) {
	sumFee, err := b.redis.Get(AccountConvertSumFee).Float64()
	if err == redis.Nil {
		return decimal.Zero, nil
	}

	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(sumFee), nil
}

func (b *box) SetAccountScanSumFee(sumFee decimal.Decimal, exp time.Duration) error {
	sumFeeF, _ := sumFee.Float64()
	err := b.redis.Set(AccountScanSumFee, sumFeeF, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) GetAccountScanSumFee() (decimal.Decimal, error) {
	sumFee, err := b.redis.Get(AccountScanSumFee).Float64()
	if err == redis.Nil {
		return decimal.Zero, nil
	}

	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromFloat(sumFee), nil
}

func (b *box) SetEthLatestHeight(height uint64, exp time.Duration) error {
	err := b.redis.Set(EthLatestHeightEth, height, exp).Err()
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

func (b *box) SetAccountGroupWinTime(height uint64, exp time.Duration) error {
	err := b.redis.Set(AccountGroupWinTime, height, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) GetAccountGroupWinTime() (time.Duration, error) {
	winTime, err := b.redis.TTL(AccountGroupWinTime).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return winTime, nil
}
