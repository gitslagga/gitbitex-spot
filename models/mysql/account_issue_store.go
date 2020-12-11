package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccountIssue(userId int64, currency string) (*models.AccountIssue, error) {
	var account models.AccountIssue
	err := s.db.Raw("SELECT * FROM g_account_issue WHERE user_id=? AND currency=?", userId,
		currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) GetAccountsIssueByUserId(userId int64) ([]*models.AccountIssue, error) {
	db := s.db.Where("user_id=?", userId).Order("id ASC")

	var accounts []*models.AccountIssue
	err := db.Find(&accounts).Error
	return accounts, err
}

func (s *Store) GetAccountIssueForUpdate(userId int64, currency string) (*models.AccountIssue, error) {
	var account models.AccountIssue
	err := s.db.Raw("SELECT * FROM g_account_issue WHERE user_id=? AND currency=? FOR UPDATE", userId, currency).Scan(&account).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &account, err
}

func (s *Store) AddAccountIssue(account *models.AccountIssue) error {
	return s.db.Create(account).Error
}

func (s *Store) UpdateAccountIssue(account *models.AccountIssue) error {
	return s.db.Save(account).Error
}
