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

func (s *Store) AddIssueLog(issueLog *models.IssueLog) error {
	return s.db.Create(issueLog).Error
}
