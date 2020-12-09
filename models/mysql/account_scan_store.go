package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetAccountScanByUserId(userId, beforeId, afterId, limit int64) ([]*models.AccountScan, error) {
	db := s.db.Where("user_id=?", userId)

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var accountScans []*models.AccountScan
	err := db.Order("id DESC").Limit(limit).Find(&accountScans).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accountScans, err
}

func (s *Store) GetAccountScanSumNumber(userId int64) (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_account_scan WHERE "+
		"created_at >= DATE_FORMAT(CURDATE(),'%Y-%m-%d') AND user_id=?", userId).Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) GetAccountScanSumFee() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(actual_number-number) as number FROM g_account_scan WHERE " +
		"created_at BETWEEN DATE_SUB(CURDATE(), INTERVAL 1 DAY) AND DATE_SUB(CURDATE(),INTERVAL 0 DAY)").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddAccountScan(accountScan *models.AccountScan) error {
	return s.db.Create(accountScan).Error
}
