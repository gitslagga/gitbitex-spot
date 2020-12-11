package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetIssueByUserId(userId, beforeId, afterId, limit int64) ([]*models.Issue, error) {
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

	var issues []*models.Issue
	err := db.Order("id DESC").Limit(limit).Find(&issues).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return issues, err
}

func (s *Store) GetIssueUsedList() ([]*models.Issue, error) {
	var issues []*models.Issue
	err := s.db.Raw("SELECT * FROM g_issue WHERE remain>0").Scan(&issues).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return issues, err
}

func (s *Store) AddIssue(issue *models.Issue) error {
	return s.db.Create(issue).Error
}

func (s *Store) UpdateIssue(issue *models.Issue) error {
	return s.db.Save(issue).Error
}

func (s *Store) GetIssueConfigs() ([]*models.IssueConfig, error) {
	var configs []*models.IssueConfig
	err := s.db.Order("id ASC").Find(&configs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return configs, err
}

func (s *Store) GetIssueLogByUserId(userId, beforeId, afterId, limit int64) ([]*models.IssueLog, error) {
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

	var issueLogs []*models.IssueLog
	err := db.Order("id DESC").Limit(limit).Find(&issueLogs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return issueLogs, err
}

func (s *Store) GetLastIssueLog(issueId int64) (*models.IssueLog, error) {
	var issueLog models.IssueLog
	err := s.db.Raw("SELECT * FROM g_issue_log WHERE issue_id=? ORDER BY id DESC LIMIT 1",
		issueId).Scan(&issueLog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &issueLog, err
}

func (s *Store) AddIssueLog(issueLog *models.IssueLog) error {
	return s.db.Create(issueLog).Error
}

func (s *Store) GetAddressHoldingByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressHolding, error) {
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

	var holdings []*models.AddressHolding
	err := db.Order("id DESC").Limit(limit).Find(&holdings).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return holdings, err
}

func (s *Store) GetLastAddressHolding() (*models.AddressHolding, error) {
	var addressHolding models.AddressHolding
	err := s.db.Raw("SELECT * FROM g_address_holding ORDER BY id DESC LIMIT 1").Scan(&addressHolding).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &addressHolding, err
}

func (s *Store) AddAddressHolding(holding *models.AddressHolding) error {
	return s.db.Create(holding).Error
}

func (s *Store) GetHoldingAccount(minHolding decimal.Decimal) ([]*models.Account, error) {
	var accounts []*models.Account
	err := s.db.Raw("SELECT * FROM g_account WHERE currency='BITE' AND available>? ORDER BY available ASC",
		minHolding.IntPart()).Scan(&accounts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accounts, err
}

func (s *Store) GetHoldingAccountAsset(minHolding decimal.Decimal) ([]*models.AccountAsset, error) {
	var accounts []*models.AccountAsset
	err := s.db.Raw("SELECT * FROM g_account_asset WHERE currency='BITE' AND available>? ORDER BY available ASC",
		minHolding.IntPart()).Scan(&accounts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accounts, err
}

func (s *Store) GetHoldingAccountPool(minHolding decimal.Decimal) ([]*models.AccountPool, error) {
	var accounts []*models.AccountPool
	err := s.db.Raw("SELECT * FROM g_account_pool WHERE currency='BITE' AND available>? ORDER BY available ASC",
		minHolding.IntPart()).Scan(&accounts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accounts, err
}

func (s *Store) GetHoldingAccountShop(minHolding decimal.Decimal) ([]*models.AccountShop, error) {
	var accounts []*models.AccountShop
	err := s.db.Raw("SELECT * FROM g_account_shop WHERE currency='BITE' AND available>? ORDER BY available ASC",
		minHolding.IntPart()).Scan(&accounts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accounts, err
}

func (s *Store) GetHoldingAccountIssue(minHolding decimal.Decimal) ([]*models.AccountIssue, error) {
	var accounts []*models.AccountIssue
	err := s.db.Raw("SELECT * FROM g_account_issue WHERE currency='BITE' AND available>? ORDER BY available ASC",
		minHolding.IntPart()).Scan(&accounts).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return accounts, err
}

func (s *Store) GetPromoteAccount() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gaa.currency, gaa.available FROM g_address ga ` +
		`INNER JOIN g_account gaa ON ga.id = gaa.user_id WHERE ga.parent_id!=0 AND gaa.currency="BITE" AND gaa.available>0`).
		Scan(&totalPower).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return totalPower, err
}

func (s *Store) GetPromoteAccountAsset() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gaa.currency, gaa.available FROM g_address ga ` +
		`INNER JOIN g_account_asset gaa ON ga.id = gaa.user_id WHERE ga.parent_id!=0 AND gaa.currency="BITE" AND gaa.available>0`).
		Scan(&totalPower).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return totalPower, err
}

func (s *Store) GetPromoteAccountPool() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gap.currency, gap.available FROM g_address ga ` +
		`INNER JOIN g_account_pool gap ON ga.id = gap.user_id WHERE ga.parent_id!=0 AND gap.currency="BITE" AND gap.available>0`).
		Scan(&totalPower).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return totalPower, err
}

func (s *Store) GetPromoteAccountShop() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gas.currency, gas.available FROM g_address ga ` +
		`INNER JOIN g_account_shop gas ON ga.id = gas.user_id WHERE ga.parent_id!=0 AND gas.currency="BITE" AND gas.available>0`).
		Scan(&totalPower).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return totalPower, err
}

func (s *Store) GetPromoteAccountIssue() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gai.currency, gai.available FROM g_address ga ` +
		`INNER JOIN g_account_issue gai ON ga.id = gai.user_id WHERE ga.parent_id!=0 AND gai.currency="BITE" AND gai.available>0`).
		Scan(&totalPower).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return totalPower, err
}

func (s *Store) GetAddressPromoteByUserId(userId, beforeId, afterId, limit int64) ([]*models.AddressPromote, error) {
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

	var promotes []*models.AddressPromote
	err := db.Order("id DESC").Limit(limit).Find(&promotes).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return promotes, err
}

func (s *Store) GetLastAddressPromote() (*models.AddressPromote, error) {
	var addressPromote models.AddressPromote
	err := s.db.Raw("SELECT * FROM g_address_promote ORDER BY id DESC LIMIT 1").Scan(&addressPromote).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &addressPromote, err
}

func (s *Store) AddAddressPromote(promote *models.AddressPromote) error {
	return s.db.Create(promote).Error
}
