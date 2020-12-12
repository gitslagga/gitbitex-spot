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

func (s *Store) GetAddressByUsername(username string) (*models.Address, error) {
	var address models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE username=?", username).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (s *Store) GetAddressByUAddressBite(addressBite string) (*models.Address, error) {
	var address models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE address_bite=?", addressBite).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &address, err
}

func (s *Store) CountAddressByMachineLevelId(machineLevelId int64) (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_address WHERE machine_level_id=?", machineLevelId).Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) GetAddressByMachineLevelId(machineLevelId int64) ([]*models.Address, error) {
	var address []*models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE machine_level_id=?", machineLevelId).Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return address, err
}

func (s *Store) CountAddressByGroupUsdt() (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_address WHERE group_usdt=1").Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) GetAddressByGroupUsdt() ([]*models.Address, error) {
	var address []*models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE group_usdt=1").Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return address, err
}

func (s *Store) CountAddressByGroupBite() (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_address WHERE group_bite=1").Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) GetAddressByGroupBite() ([]*models.Address, error) {
	var address []*models.Address
	err := s.db.Raw("SELECT * FROM g_address WHERE group_bite=1").Scan(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return address, err
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
