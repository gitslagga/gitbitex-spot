package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetMachineConvertByUserId(userId, beforeId, afterId, limit int64) ([]*models.MachineConvert, error) {
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

	var machineConverts []*models.MachineConvert
	err := db.Order("id DESC").Limit(limit).Find(&machineConverts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return machineConverts, err
}

func (s *Store) GetMachineConvertSumNum() (decimal.Decimal, error) {
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
		"created_at >= DATE_SUB(CURDATE(),INTERVAL -1 DAY)").Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddMachineConvert(machineConvert *models.MachineConvert) error {
	return s.db.Create(machineConvert).Error
}
