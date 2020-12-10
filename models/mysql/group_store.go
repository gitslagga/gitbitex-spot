package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

func (s *Store) GetGroupById(groupId int64) (*models.Group, error) {
	var group models.Group
	err := s.db.Raw("SELECT * FROM g_group WHERE id=?", groupId).Scan(&group).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (s *Store) GetGroupByUserIdCoin(userId int64, coin string) (*models.Group, error) {
	var group models.Group
	err := s.db.Raw("SELECT * FROM g_group WHERE user_id=? AND coin=? AND status=1", userId, coin).Scan(&group).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &group, err
}

func (s *Store) GetGroupByCoin(coin string, beforeId, afterId, limit int64) ([]*models.Group, error) {
	db := s.db.Where("coin=?", coin)

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var groups []*models.Group
	err := db.Order("id DESC").Limit(limit).Find(&groups).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return groups, err
}

func (s *Store) AddGroup(group *models.Group) error {
	return s.db.Create(group).Error
}

func (s *Store) UpdateGroup(group *models.Group) error {
	return s.db.Save(group).Error
}

// GroupLog
func (s *Store) GetGroupLogByUserId(userId, beforeId, afterId, limit int64) ([]*models.GroupLog, error) {
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

	var groupLogs []*models.GroupLog
	err := db.Order("id DESC").Limit(limit).Find(&groupLogs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return groupLogs, err
}

func (s *Store) GetGroupLogPublicity(beforeId, afterId, limit int64) ([]*models.GroupLog, error) {
	db := s.db.Where("status=1")

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}
	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}
	if limit <= 0 {
		limit = 10
	}

	var groupLogs []*models.GroupLog
	err := db.Order("id DESC").Limit(limit).Find(&groupLogs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return groupLogs, err
}

func (s *Store) GetGroupLogByGroupId(groupId int64) ([]*models.GroupLog, error) {
	var groupLogs []*models.GroupLog
	err := s.db.Raw("SELECT * FROM g_group_log WHERE group_id=?", groupId).Scan(&groupLogs).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return groupLogs, err
}

func (s *Store) GetGroupLogSumNum(coin string) (decimal.Decimal, error) {
	var number models.SumNumber
	err := s.db.Raw("SELECT SUM(number) as number FROM g_group_log WHERE "+
		"created_at BETWEEN DATE_SUB(CURDATE(), INTERVAL 1 DAY) AND DATE_SUB(CURDATE(),INTERVAL 0 DAY) AND coin=?", coin).Scan(&number).Error
	if err == gorm.ErrRecordNotFound {
		return decimal.Zero, nil
	}

	return number.Number, err
}

func (s *Store) AddGroupLog(groupLog *models.GroupLog) error {
	return s.db.Create(groupLog).Error
}

func (s *Store) UpdateGroupLog(groupLog *models.GroupLog) error {
	return s.db.Save(groupLog).Error
}
