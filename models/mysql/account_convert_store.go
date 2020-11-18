package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetAccountConvertByUserId(userId int64) ([]*models.AccountConvert, error) {
	db := s.db.Where("user_id=?", userId).Order("id DESC")

	var accountConvert []*models.AccountConvert
	err := db.Find(&accountConvert).Error
	return accountConvert, err
}

func (s *Store) GetAccountConvertSumNumber() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_account_convert WHERE " +
		"DATE_FORMAT(created_at,'%Y-%m-%d') = DATE_FORMAT(CURDATE(),'%Y-%m-%d')").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddAccountConvert(accountConvert *models.AccountConvert) error {
	return s.db.Create(accountConvert).Error
}
