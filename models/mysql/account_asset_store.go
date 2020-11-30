package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
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

func (s *Store) GetIssueAccountAsset() ([]*models.AccountAsset, error) {
	var assets []*models.AccountAsset
	err := s.db.Raw("SELECT * FROM g_account_asset WHERE currency='USDT' AND available>0").Scan(&assets).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return assets, err
}

func (s *Store) SumIssueAccountAsset() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(available) FROM g_account_asset WHERE currency='USDT' AND available>0").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}
	return number.Number, err
}

func (s *Store) AddAccountAsset(account *models.AccountAsset) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccountAsset(account *models.AccountAsset) error {
	return s.db.Save(account).Error
}
