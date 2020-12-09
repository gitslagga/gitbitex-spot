package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAddressReleaseByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressRelease, error) {
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

	var releases []*models.AddressRelease
	err := db.Order("id DESC").Limit(limit).Find(&releases).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return releases, err
}

func (s *Store) GetLastAddressRelease(releaseType int) (*models.AddressRelease, error) {
	var release models.AddressRelease
	err := s.db.Raw("SELECT * FROM g_address_release WHERE type=? ORDER BY id DESC LIMIT 1",
		releaseType).Scan(&release).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &release, err
}

func (s *Store) AddAddressRelease(release *models.AddressRelease) error {
	return s.db.Create(release).Error
}
