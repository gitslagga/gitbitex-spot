package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineAddressByUserId(userId int64) ([]*models.MachineAddress, error) {
	var machineAddress []*models.MachineAddress
	err := s.db.Raw("SELECT * FROM g_machine_address WHERE user_id=?", userId).Scan(&machineAddress).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineAddress, err
}

func (s *Store) GetMachineAddressUsedCount(userId int64) (int, error) {
	var count models.TotalCount
	err := s.db.Raw("SELECT COUNT(*) as count FROM g_machine_address WHERE user_id=? AND day!=0", userId).Scan(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count.Count, err
}

func (s *Store) AddMachineAddress(machineAddress *models.MachineAddress) error {
	return s.db.Create(machineAddress).Error
}

func (s *Store) UpdateMachineAddress(machineAddress *models.MachineAddress) error {
	return s.db.Save(machineAddress).Error

}