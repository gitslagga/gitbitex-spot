package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccountShop(userId int64, currency string) (*models.AccountShop, error) {
	var account models.AccountShop
	err := s.db.Raw("SELECT * FROM g_account_shop WHERE user_id=? AND currency=?", userId,
		currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) GetAccountsShopByUserId(userId int64) ([]*models.AccountShop, error) {
	db := s.db.Where("user_id=?", userId).Order("id ASC")

	var accounts []*models.AccountShop
	err := db.Find(&accounts).Error
	return accounts, err
}

func (s *Store) GetAccountShopForUpdate(userId int64, currency string) (*models.AccountShop, error) {
	var account models.AccountShop
	err := s.db.Raw("SELECT * FROM g_account_shop WHERE user_id=? AND currency=? FOR UPDATE", userId, currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) AddAccountShop(account *models.AccountShop) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccountShop(account *models.AccountShop) error {
	return s.db.Save(account).Error
}
