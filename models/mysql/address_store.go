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

func (s *Store) AddAddress(address *models.Address) error {
	return s.db.Create(address).Error
}

func (s *Store) UpdateAddress(address *models.Address) error {
	return s.db.Save(address).Error

}
