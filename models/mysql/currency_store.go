package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetValidAddressConfig() ([]*models.AddressConfig, error) {
	var configs []*models.AddressConfig
	err := s.db.Find(&configs).Where("status=2").Order("id ASC").Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return configs, err
}

func (s *Store) GetAddressConfigByCoin(coin string) (*models.AddressConfig, error) {
	var config models.AddressConfig
	err := s.db.Where("coin=?", coin).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (s *Store) GetAddressConfigByContract(contract string) (*models.AddressConfig, error) {
	var config models.AddressConfig
	err := s.db.Where("contract_address=?", contract).First(&config).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &config, err
}

func (s *Store) UpdateAddressConfig(config *models.AddressConfig) error {
	return s.db.Save(config).Error
}

func (s *Store) AddAddressCollect(collect *models.AddressCollect) error {
	return s.db.Create(collect).Error
}

func (s *Store) GetAddressDepositsByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressDeposit, error) {
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

	var deposits []*models.AddressDeposit
	err := db.Order("id DESC").Limit(limit).Find(&deposits).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return deposits, err
}

func (s *Store) GetAddressDepositsByBNStatus(blockNum uint64, status int) ([]*models.AddressDeposit, error) {
	var deposits []*models.AddressDeposit
	err := s.db.Where("block_num=? AND status=?", blockNum, status).Find(&deposits).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return deposits, err
}

func (s *Store) AddAddressDeposit(deposit *models.AddressDeposit) error {
	return s.db.Create(deposit).Error
}

func (s *Store) UpdateAddressDeposit(deposit *models.AddressDeposit) error {
	return s.db.Save(deposit).Error
}

func (s *Store) GetAddressWithdrawsByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressWithdraw, error) {
	db := s.db.Where("user_id=?", userId)

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var withdraws []*models.AddressWithdraw
	err := db.Order("id DESC").Limit(limit).Find(&withdraws).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return withdraws, err
}

func (s *Store) GetAddressWithdrawsByOrderSN(orderSN string) (*models.AddressWithdraw, error) {
	var withdraw *models.AddressWithdraw
	err := s.db.Where("order_sn =?", orderSN).Find(&withdraw).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return withdraw, err
}

func (s *Store) AddAddressWithdraw(withdraw *models.AddressWithdraw) error {
	return s.db.Create(withdraw).Error
}

func (s *Store) UpdateAddressWithdraw(withdraw *models.AddressWithdraw) error {
	return s.db.Save(withdraw).Error
}
