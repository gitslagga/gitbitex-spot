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
	err := b.redis.Set(AccountConvertSumFee, sumFee, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

func (b *box) GetAccountConvertSumFee() (decimal.Decimal, error) {
	sumFee, err := b.redis.Get(AccountConvertSumFee).Float64()
	if err != nil {
		return decimal.Zero, err
	}

	if err == redis.Nil {
		return decimal.Zero, nil
	}

	return decimal.NewFromFloat(sumFee), nil
}
