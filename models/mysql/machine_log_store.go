package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineLogByUserId(userId, beforeId, afterId, limit int64) ([]*models.MachineLog, error) {
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

	var machineLogs []*models.MachineLog
	err := db.Order("id DESC").Limit(limit).Find(&machineLogs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineLogs, err
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
