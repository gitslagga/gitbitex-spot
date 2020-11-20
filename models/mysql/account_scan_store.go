package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetAccountScanByUserId(userId int64) ([]*models.AccountScan, error) {
	db := s.db.Where("user_id=?", userId).Order("id DESC")

	var accountScan []*models.AccountScan
	err := db.Find(&accountScan).Error
	return accountScan, err
}

func (s *Store) GetAccountScanSumNumber(userId int64) (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_machine_scan WHERE "+
		"DATE_FORMAT(created_at,'%Y-%m-%d') = DATE_FORMAT(CURDATE(),'%Y-%m-%d') AND user_id=?", userId).Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) GetAccountScanSumFee() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(actual-number) as number FROM g_machine_scan WHERE " +
		"DATE_FORMAT(created_at,'%Y-%m-%d') = DATE_FORMAT(CURDATE(),'%Y-%m-%d')").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddAccountScan(accountScan *models.AccountScan) error {
	return s.db.Create(accountScan).Error
}
