package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccount(userId int64, currency string) (*models.Account, error) {
	var account models.Account
	err := s.db.Raw("SELECT * FROM g_account WHERE user_id=? AND currency=?", userId,
		currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) GetAccountsByUserId(userId int64) ([]*models.Account, error) {
	db := s.db.Where("user_id=?", userId).Order("id ASC")

	var accounts []*models.Account
	err := db.Find(&accounts).Error
	return accounts, err
}

func (s *Store) GetAccountForUpdate(userId int64, currency string) (*models.Account, error) {
	var account models.Account
	err := s.db.Raw("SELECT * FROM g_account WHERE user_id=? AND currency=? FOR UPDATE", userId, currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) AddAccount(account *models.Account) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccount(account *models.Account) error {
	return s.db.Save(account).Error
}
