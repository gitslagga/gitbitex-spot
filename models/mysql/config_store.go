package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetConfigs() ([]*models.Config, error) {
	var configs []*models.Config
	err := s.db.Find(&configs).Order("id ASC").Error
	return configs, err
}

func (s *Store) GetConfigById(id int64) (*models.Config, error) {
	var config models.Config
	err := s.db.Raw("SELECT * FROM g_config WHERE id=?", id).Scan(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (s *Store) UpdateConfig(config *models.Config) error {
	return s.db.Save(config).Error
}
