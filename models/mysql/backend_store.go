package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
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

func (s *Store) GetTotalPowerList() ([]*models.TotalPower, error) {
	var totalPower []*models.TotalPower
	err := s.db.Raw(`SELECT ga.id,ga.parent_id,ga.parent_ids, gaa.currency, gaa.available FROM g_address ga ` +
		`INNER JOIN g_account_asset gaa ON ga.id = gaa.user_id WHERE ga.parent_id!=0 AND gaa.currency="BITE" AND gaa.available>0`).
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
