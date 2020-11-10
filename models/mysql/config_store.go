package mysql

import "github.com/gitslagga/gitbitex-spot/models"

func (s *Store) GetConfigs() ([]*models.Config, error) {
	var configs []*models.Config
	err := s.db.Find(&configs).Error
	return configs, err
}

func (s *Store) UpdateConfig(config *models.Config) error {
	return s.db.Save(config).Error
}
