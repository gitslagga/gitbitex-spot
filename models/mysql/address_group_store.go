package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAddressGroupByUserId(userId int64) ([]*models.AddressGroup, error) {
	db := s.db.Where("user_id=?", userId).Order("id ASC")

	var groups []*models.AddressGroup
	err := db.Find(&groups).Error
	return groups, err
}

func (s *Store) GetAddressGroupForUpdate(currency string) (*models.AddressGroup, error) {
	var group models.AddressGroup
	err := s.db.Raw("SELECT * FROM g_address_group WHERE currency=? FOR UPDATE ORDER BY id DESC LIMIT 1",
		currency).Scan(&group).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (s *Store) AddAddressGroup(group *models.AddressGroup) error {
	return s.db.Create(group).Error
}

func (s *Store) UpdateAddressGroup(group *models.AddressGroup) error {
	return s.db.Save(group).Error
}
