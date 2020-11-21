package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAccountTransferByUserId(userId, beforeId, afterId, limit int64) ([]*models.AccountTransfer, error) {

	db := s.db.Where("user_id =?", userId)

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var accountTransfers []*models.AccountTransfer
	err := db.Order("id DESC").Limit(limit).Find(&accountTransfers).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accountTransfers, err
}

func (s *Store) AddAccountTransfer(accountTransfer *models.AccountTransfer) error {
	return s.db.Create(accountTransfer).Error
}
