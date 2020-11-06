package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccountPool(userId int64, currency string) (*models.AccountPool, error) {
	var account models.AccountPool
	err := s.db.Raw("SELECT * FROM g_account_pool WHERE user_id=? AND currency=?", userId,
		currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) GetAccountsPoolByUserId(userId int64) ([]*models.AccountPool, error) {
	db := s.db.Where("user_id=?", userId)

	var accounts []*models.AccountPool
	err := db.Find(&accounts).Error
	return accounts, err
}

func (s *Store) GetAccountPoolForUpdate(userId int64, currency string) (*models.AccountPool, error) {
	var account models.AccountPool
	err := s.db.Raw("SELECT * FROM g_account_pool WHERE user_id=? AND currency=? FOR UPDATE", userId, currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) AddAccountPool(account *models.AccountPool) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccountPool(account *models.AccountPool) error {
	return s.db.Save(account).Error
}
