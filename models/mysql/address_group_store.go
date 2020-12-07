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
	err := s.db.Raw("SELECT * FROM g_address_group WHERE coin=? ORDER BY id DESC LIMIT 1 FOR UPDATE",
		currency).Scan(&group).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (s *Store) GetAddressGroupByUserIdOrderSN(userId int64, orderSN string) (*models.AddressGroup, error) {
	var group models.AddressGroup
	err := s.db.Raw("SELECT * FROM g_address_group WHERE user_id=? AND order_sn=?", userId, orderSN).Scan(&group).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (s *Store) GetAddressGroupsByOrderSN(orderSN string) ([]*models.AddressGroup, error) {
	var groups []*models.AddressGroup
	err := s.db.Raw("SELECT * FROM g_address_group WHERE order_sn=?", orderSN).Scan(&groups).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return groups, err
}

func (s *Store) AddAddressGroup(group *models.AddressGroup) error {
	return s.db.Create(group).Error
}

func (s *Store) UpdateAddressGroup(group *models.AddressGroup) error {
	return s.db.Save(group).Error
}
