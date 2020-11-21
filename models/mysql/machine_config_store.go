package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetMachineConfigs() ([]*models.MachineConfig, error) {
	var configs []*models.MachineConfig
	err := s.db.Find(&configs).Order("id ASC").Error
	return configs, err
}

func (s *Store) GetMachineConfigById(id int64) (*models.MachineConfig, error) {
	var config models.MachineConfig
	err := s.db.Raw("SELECT * FROM g_machine_config WHERE id=?", id).Scan(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (s *Store) UpdateMachineConfig(config *models.MachineConfig) error {
	return s.db.Save(config).Error
}
