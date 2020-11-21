package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetValidCurrencies() ([]*models.Currency, error) {
	var currencies []*models.Currency
	err := s.db.Find(&currencies).Where("status=2").Order("id ASC").Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return currencies, err
}

func (s *Store) GetCurrencyByCoin(coin string) (*models.Currency, error) {
	var currency models.Currency
	err := s.db.Raw("SELECT * FROM g_currency WHERE coin=?", coin).Scan(&currency).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &currency, err
}

func (s *Store) UpdateCurrency(currency *models.Currency) error {
	return s.db.Save(currency).Error
}

func (s *Store) AddCurrencyCollect(currencyCollect *models.CurrencyCollect) error {
	return s.db.Create(currencyCollect).Error
}

func (s *Store) GetCurrencyDepositsByUserId(userId int64) ([]*models.CurrencyDeposit, error) {
	var currencyDeposits []*models.CurrencyDeposit
	err := s.db.Raw("SELECT * FROM g_currency_deposit WHERE user_id=?", userId).Scan(&currencyDeposits).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return currencyDeposits, err
}

func (s *Store) AddCurrencyDeposit(currencyDeposit *models.CurrencyDeposit) error {
	return s.db.Create(currencyDeposit).Error
}

func (s *Store) UpdateCurrencyDeposit(currencyDeposit *models.CurrencyDeposit) error {
	return s.db.Save(currencyDeposit).Error
}

func (s *Store) GetCurrencyWithdrawsByUserId(userId int64) ([]*models.CurrencyWithdraw, error) {
	var currencyWithdraws []*models.CurrencyWithdraw
	err := s.db.Raw("SELECT * FROM g_currency_withdraw WHERE user_id=?", userId).Scan(&currencyWithdraws).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return currencyWithdraws, err
}

func (s *Store) AddCurrencyWithdraw(currencyWithdraw *models.CurrencyWithdraw) error {
	return s.db.Create(currencyWithdraw).Error
}

func (s *Store) UpdateCurrencyWithdraw(currencyWithdraw *models.CurrencyWithdraw) error {
	return s.db.Save(currencyWithdraw).Error
}
