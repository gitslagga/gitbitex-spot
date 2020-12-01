package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAddressListByAddress(address string) (*models.AddressList, error) {
	var addressList models.AddressList
	err := s.db.Raw("SELECT * FROM g_address_list WHERE address=?", address).Scan(&addressList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &addressList, err
}

func (s *Store) GetAddressListById(id int64) (*models.AddressList, error) {
	var addressList models.AddressList
	err := s.db.Raw("SELECT * FROM g_address_list WHERE id=?", id).Scan(&addressList).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &addressList, err
}

func (s *Store) GetAddressListByUserId(userId int64) ([]*models.AddressList, error) {
	var addressLists []*models.AddressList
	err := s.db.Raw("SELECT * FROM g_address_list WHERE user_id=?", userId).Scan(&addressLists).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return addressLists, err
}

func (s *Store) CountAddressListByUserId(userId int64) (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_address_list WHERE user_id=?", userId).Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) AddAddressList(addressList *models.AddressList) error {
	return s.db.Create(addressList).Error
}

func (s *Store) UpdateAddressList(addressList *models.AddressList) error {
	return s.db.Save(addressList).Error
}

func (s *Store) DeleteAddressList(addressList *models.AddressList) error {
	return s.db.Where("id=", addressList.Id).Delete(addressList).Error
}
