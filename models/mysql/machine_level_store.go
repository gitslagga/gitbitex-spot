package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineLevel() ([]*models.MachineLevel, error) {
	var machineLevel []*models.MachineLevel
	err := s.db.Raw("SELECT * FROM g_machine_level").Order("id ASC").Scan(&machineLevel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineLevel, err
}

func (s *Store) GetMachineLevelById(machineLevelId int64) (*models.MachineLevel, error) {
	var machineLevel models.MachineLevel
	err := s.db.Raw("SELECT * FROM g_machine_level WHERE id=?", machineLevelId).Scan(&machineLevel).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &machineLevel, err
}
