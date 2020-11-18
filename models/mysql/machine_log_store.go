package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineLogByUserId(userId int64) ([]*models.MachineLog, error) {
	var machineLog []*models.MachineLog
	err := s.db.Raw("SELECT * FROM g_machine_log WHERE user_id=?", userId).Order("id DESC").Scan(&machineLog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineLog, err
}

func (s *Store) GetLastMachineLog(machineAddressId int64) (*models.MachineLog, error) {
	var machineLog models.MachineLog
	err := s.db.Raw("SELECT * FROM g_machine_log WHERE machine_address_id=? ORDER BY id DESC LIMIT 1",
		machineAddressId).Scan(&machineLog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &machineLog, err
}

func (s *Store) AddMachineLog(machineLog *models.MachineLog) error {
	return s.db.Create(machineLog).Error
}
