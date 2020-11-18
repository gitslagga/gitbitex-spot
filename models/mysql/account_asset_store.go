package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccountAsset(userId int64, currency string) (*models.AccountAsset, error) {
	var account models.AccountAsset
	err := s.db.Raw("SELECT * FROM g_account_asset WHERE user_id=? AND currency=?", userId,
		currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) GetAccountsAssetByUserId(userId int64) ([]*models.AccountAsset, error) {
	db := s.db.Where("user_id=?", userId).Order("id ASC")

	var accounts []*models.AccountAsset
	err := db.Find(&accounts).Error
	return accounts, err
}

func (s *Store) GetAccountAssetForUpdate(userId int64, currency string) (*models.AccountAsset, error) {
	var account models.AccountAsset
	err := s.db.Raw("SELECT * FROM g_account_asset WHERE user_id=? AND currency=? FOR UPDATE", userId, currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) AddAccountAsset(account *models.AccountAsset) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccountAsset(account *models.AccountAsset) error {
	return s.db.Save(account).Error
}
