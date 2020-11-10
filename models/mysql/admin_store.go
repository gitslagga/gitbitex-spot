package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetAdmin(username string) (*models.Admin, error) {
	var admin models.Admin
	err := s.db.Raw("SELECT * FROM g_admin WHERE username=?", username).Scan(&admin).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &admin, err
}

func (s *Store) UpdateAdmin(admin *models.Admin) error {
	return s.db.Save(admin).Error
}
