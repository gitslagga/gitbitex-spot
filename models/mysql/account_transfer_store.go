package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
)

func (s *Store) GetAccountTransferByUserId(userId int64) ([]*models.AccountTransfer, error) {
	db := s.db.Where("user_id=?", userId).Order("id DESC")

	var accountTransfer []*models.AccountTransfer
	err := db.Find(&accountTransfer).Error
	return accountTransfer, err
}

func (s *Store) AddAccountTransfer(accountTransfer *models.AccountTransfer) error {
	return s.db.Create(accountTransfer).Error
}
