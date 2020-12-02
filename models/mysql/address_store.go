package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAddressByAddress(addr string) (*models.Address, error) {
	var address models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE address=?", addr).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (s *Store) GetAddressById(id int64) (*models.Address, error) {
	var address models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE id=?", id).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (s *Store) CountAddressByMachineLevelId(machineLevelId int64) (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_address WHERE g_machine_level_id=?", machineLevelId).Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) GetAddressByParentId(parentId int64) ([]*models.Address, error) {
	var address []*models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE parent_id=?", parentId).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return address, err
}

func (s *Store) AddAddress(address *models.Address) error {
	return s.db.Create(address).Error
}

func (s *Store) UpdateAddress(address *models.Address) error {
	return s.db.Save(address).Error

}

func (s *Store) GetAddressHoldingByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressHolding, error) {
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

	var holdings []*models.AddressHolding
	err := db.Order("id DESC").Limit(limit).Find(&holdings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return holdings, err
}

func (s *Store) AddAddressHolding(holding *models.AddressHolding) error {
	return s.db.Create(holding).Error
}

func (s *Store) GetAddressPromoteByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressPromote, error) {
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

	var promotes []*models.AddressPromote
	err := db.Order("id DESC").Limit(limit).Find(&promotes).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return promotes, err
}

func (s *Store) AddAddressPromote(promote *models.AddressPromote) error {
	return s.db.Create(promote).Error
}
