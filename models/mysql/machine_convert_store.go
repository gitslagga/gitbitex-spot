package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetMachineConvertByUserId(userId int64) ([]*models.MachineConvert, error) {
	db := s.db.Where("user_id=?", userId).Order("id DESC")

	var machineConvert []*models.MachineConvert
	err := db.Find(&machineConvert).Error
	return machineConvert, err
}

func (s *Store) GetMachineConvertSumNumber() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_machine_convert WHERE " +
		"DATE_FORMAT(created_at,'%Y-%m-%d') = DATE_FORMAT(CURDATE(),'%Y-%m-%d')").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) GetMachineConvertSumFee() (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(amount-number) as number FROM g_machine_convert WHERE " +
		"DATE_FORMAT(created_at,'%Y-%m-%d') = DATE_FORMAT(CURDATE(),'%Y-%m-%d')").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddMachineConvert(machineConvert *models.MachineConvert) error {
	return s.db.Create(machineConvert).Error
}
