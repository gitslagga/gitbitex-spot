package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetBuyMachine() ([]*models.Machine, error) {
	var machines []*models.Machine
	err := s.db.Raw("SELECT * FROM g_machine WHERE buy_quantity>0").Order("id ASC").Scan(&machines).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machines, err
}

func (s *Store) GetMachineById(machineId int64) (*models.Machine, error) {
	var machine models.Machine
	err := s.db.Raw("SELECT * FROM g_machine WHERE id=?", machineId).Scan(&machine).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &machine, err
}
