package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetAddressGroupByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressGroup, error) {
	db := s.db.Where("user_id=?", userId)

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var groups []*models.AddressGroup
	err := db.Order("id DESC").Limit(limit).Find(&groups).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
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

func (s *Store) GetAddressGroupSumNum(coin string) (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_address_group WHERE "+
		"created_at >= DATE_SUB(CURDATE(),INTERVAL -1 DAY) AND coin=?", coin).Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddAddressGroup(group *models.AddressGroup) error {
	return s.db.Create(group).Error
}

func (s *Store) UpdateAddressGroup(group *models.AddressGroup) error {
	return s.db.Save(group).Error
}
