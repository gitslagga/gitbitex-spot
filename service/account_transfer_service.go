package service

import (
	"errors"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/shopspring/decimal"
)

func GetAccountTransferByUserId(userId, before, after, limit int64) ([]*models.AccountTransfer, error) {
	return mysql.SharedStore().GetAccountTransferByUserId(userId, before, after, limit)
}

func AddAccountTransfer(accountTransfer *models.AccountTransfer) error {
	return mysql.SharedStore().AddAccountTransfer(accountTransfer)
}

func AccountTransfer(userId int64, from, to int, currency string, amount float64) error {
	var err error
	number := decimal.NewFromFloat(amount)

	switch from {
	case models.AccountAssetTransfer:
		if to == models.AccountPoolTransfer {
			err = transferFromAssetToPool(userId, from, to, currency, number)
		} else if to == models.AccountSpotTransfer {
			err = transferFromAssetToSpot(userId, from, to, currency, number)
		} else if to == models.AccountShopTransfer {
			err = transferFromAssetToShop(userId, from, to, currency, number)
		}
	case models.AccountPoolTransfer:
		if to == models.AccountAssetTransfer {
			err = transferFromPoolToAsset(userId, from, to, currency, number)
		} else if to == models.AccountSpotTransfer {
			err = transferFromPoolToSpot(userId, from, to, currency, number)
		} else if to == models.AccountShopTransfer {
			err = transferFromPoolToShop(userId, from, to, currency, number)
		}
	case models.AccountSpotTransfer:
		if to == models.AccountAssetTransfer {
			err = transferFromSpotToAsset(userId, from, to, currency, number)
		} else if to == models.AccountPoolTransfer {
			err = transferFromSpotToPool(userId, from, to, currency, number)
		} else if to == models.AccountShopTransfer {
			err = transferFromSpotToShop(userId, from, to, currency, number)
		}
		//case models.AccountShopTransfer:
		//	if to == models.AccountAssetTransfer {
		//		err = transferFromShopToAsset(userId, from, to, currency, number)
		//	} else if to == models.AccountPoolTransfer {
		//		err = transferFromShopToPool(userId, from, to, currency, number)
		//	} else if to == models.AccountSpotTransfer {
		//		err = transferFromShopToSpot(userId, from, to, currency, number)
		//	}
	}

	return err
}

func transferFromAssetToPool(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	asset, err := db.GetAccountAssetForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if asset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	asset.Available = asset.Available.Sub(number)
	err = db.UpdateAccountAsset(asset)
	if err != nil {
		return err
	}

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	pool.Available = pool.Available.Add(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromAssetToSpot(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	asset, err := db.GetAccountAssetForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if asset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	asset.Available = asset.Available.Sub(number)
	err = db.UpdateAccountAsset(asset)
	if err != nil {
		return err
	}

	spot, err := db.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	spot.Available = spot.Available.Add(number)
	err = db.UpdateAccount(spot)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromAssetToShop(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	asset, err := db.GetAccountAssetForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if asset.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	asset.Available = asset.Available.Sub(number)
	err = db.UpdateAccountAsset(asset)
	if err != nil {
		return err
	}

	shop, err := db.GetAccountShopForUpdate(userId, currency)
	if err != nil {
		return err
	}
	shop.Available = shop.Available.Add(number)
	err = db.UpdateAccountShop(shop)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromPoolToAsset(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if pool.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	pool.Available = pool.Available.Sub(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	asset, err := db.GetAccountAssetForUpdate(userId, currency)
	if err != nil {
		return err
	}
	asset.Available = asset.Available.Add(number)
	err = db.UpdateAccountAsset(asset)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromPoolToSpot(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if pool.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	pool.Available = pool.Available.Sub(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	spot, err := db.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	spot.Available = spot.Available.Add(number)
	err = db.UpdateAccount(spot)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromPoolToShop(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if pool.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	pool.Available = pool.Available.Sub(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	shop, err := db.GetAccountShopForUpdate(userId, currency)
	if err != nil {
		return err
	}
	shop.Available = shop.Available.Add(number)
	err = db.UpdateAccountShop(shop)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromSpotToAsset(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if pool.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	pool.Available = pool.Available.Sub(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	asset, err := db.GetAccountAssetForUpdate(userId, currency)
	if err != nil {
		return err
	}
	asset.Available = asset.Available.Add(number)
	err = db.UpdateAccountAsset(asset)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromSpotToPool(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	spot, err := db.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if spot.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	spot.Available = spot.Available.Sub(number)
	err = db.UpdateAccount(spot)
	if err != nil {
		return err
	}

	pool, err := db.GetAccountPoolForUpdate(userId, currency)
	if err != nil {
		return err
	}
	pool.Available = pool.Available.Add(number)
	err = db.UpdateAccountPool(pool)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

func transferFromSpotToShop(userId int64, from, to int, currency string, number decimal.Decimal) error {
	db, err := mysql.SharedStore().BeginTx()
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	spot, err := db.GetAccountForUpdate(userId, currency)
	if err != nil {
		return err
	}
	if spot.Available.LessThan(number) {
		return errors.New("资产余额不足|Insufficient number of asset")
	}
	spot.Available = spot.Available.Sub(number)
	err = db.UpdateAccount(spot)
	if err != nil {
		return err
	}

	shop, err := db.GetAccountShopForUpdate(userId, currency)
	if err != nil {
		return err
	}
	shop.Available = shop.Available.Add(number)
	err = db.UpdateAccountShop(shop)
	if err != nil {
		return err
	}

	err = db.AddAccountTransfer(&models.AccountTransfer{
		UserId:   userId,
		From:     from,
		To:       to,
		Currency: currency,
		Number:   number,
	})
	if err != nil {
		return err
	}

	return db.CommitTx()
}

//func transferFromShopToAsset(userId int64, from, to int, currency string, number decimal.Decimal) error {
//	db, err := mysql.SharedStore().BeginTx()
//	if err != nil {
//		return err
//	}
//	defer func() { _ = db.Rollback() }()
//
//	shop, err := db.GetAccountShopForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	if shop.Available.LessThan(number) {
//		return errors.New("资产余额不足|Insufficient number of asset")
//	}
//	shop.Available = shop.Available.Sub(number)
//	err = db.UpdateAccountShop(shop)
//	if err != nil {
//		return err
//	}
//
//	asset, err := db.GetAccountAssetForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	asset.Available = asset.Available.Add(number)
//	err = db.UpdateAccountAsset(asset)
//	if err != nil {
//		return err
//	}
//
//	err = db.AddAccountTransfer(&models.AccountTransfer{
//		UserId:   userId,
//		From:     from,
//		To:       to,
//		Currency: currency,
//		Number:   number,
//	})
//	if err != nil {
//		return err
//	}
//
//	return db.CommitTx()
//}
//
//func transferFromShopToPool(userId int64, from, to int, currency string, number decimal.Decimal) error {
//	db, err := mysql.SharedStore().BeginTx()
//	if err != nil {
//		return err
//	}
//	defer func() { _ = db.Rollback() }()
//
//	shop, err := db.GetAccountShopForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	if shop.Available.LessThan(number) {
//		return errors.New("资产余额不足|Insufficient number of asset")
//	}
//	shop.Available = shop.Available.Sub(number)
//	err = db.UpdateAccountShop(shop)
//	if err != nil {
//		return err
//	}
//
//	pool, err := db.GetAccountPoolForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	pool.Available = pool.Available.Add(number)
//	err = db.UpdateAccountPool(pool)
//	if err != nil {
//		return err
//	}
//
//	err = db.AddAccountTransfer(&models.AccountTransfer{
//		UserId:   userId,
//		From:     from,
//		To:       to,
//		Currency: currency,
//		Number:   number,
//	})
//	if err != nil {
//		return err
//	}
//
//	return db.CommitTx()
//}
//
//func transferFromShopToSpot(userId int64, from, to int, currency string, number decimal.Decimal) error {
//	db, err := mysql.SharedStore().BeginTx()
//	if err != nil {
//		return err
//	}
//	defer func() { _ = db.Rollback() }()
//
//	shop, err := db.GetAccountShopForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	if shop.Available.LessThan(number) {
//		return errors.New("资产余额不足|Insufficient number of asset")
//	}
//	shop.Available = shop.Available.Sub(number)
//	err = db.UpdateAccountShop(shop)
//	if err != nil {
//		return err
//	}
//
//	spot, err := db.GetAccountForUpdate(userId, currency)
//	if err != nil {
//		return err
//	}
//	spot.Available = spot.Available.Add(number)
//	err = db.UpdateAccount(spot)
//	if err != nil {
//		return err
//	}
//
//	err = db.AddAccountTransfer(&models.AccountTransfer{
//		UserId:   userId,
//		From:     from,
//		To:       to,
//		Currency: currency,
//		Number:   number,
//	})
//	if err != nil {
//		return err
//	}
//
//	return db.CommitTx()
//}
