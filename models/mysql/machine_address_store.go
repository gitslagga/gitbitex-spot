package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineAddressByUserId(userId, beforeId, afterId, limit int64) ([]*models.MachineAddress, error) {
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

	var machineAddresses []*models.MachineAddress
	err := db.Order("id DESC").Limit(limit).Find(&machineAddresses).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineAddresses, err
}

func (s *Store) CountMachineAddressUsed(userId int64, machineId int64) (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_machine_address WHERE user_id=? AND machine_id=? AND is_buy=1",
		userId, machineId).Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) GetMachineAddressUsedList() ([]*models.MachineAddress, error) {
	var machineAddress []*models.MachineAddress
	err := s.db.Raw("SELECT * FROM g_machine_address WHERE day>0").Scan(&machineAddress).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineAddress, err
}

func (s *Store) AddMachineAddress(machineAddress *models.MachineAddress) error {
	return s.db.Create(machineAddress).Error
}

func (s *Store) UpdateMachineAddress(machineAddress *models.MachineAddress) error {
	return s.db.Save(machineAddress).Error

}
