package mysql

import (
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
)

func (s *Store) GetOrderById(orderId int64) (*models.Order, error) {
	var order models.Order
	err := s.db.Raw("SELECT * FROM g_order WHERE id=?", orderId).Scan(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (s *Store) GetOrderByClientOid(userId int64, clientOid string) (*models.Order, error) {
	var order models.Order
	err := s.db.Raw("SELECT * FROM g_order WHERE user_id=? AND client_oid=?", userId, clientOid).Scan(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (s *Store) GetOrderByIdForUpdate(orderId int64) (*models.Order, error) {
	var order models.Order
	err := s.db.Raw("SELECT * FROM g_order WHERE id=? FOR UPDATE", orderId).Scan(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (s *Store) GetOrdersByUserId(userId int64, statuses []models.OrderStatus, side *models.Side, productId string,
	beforeId, afterId int64, limit int) ([]*models.Order, error) {
	db := s.db.Where("user_id =?", userId)

	if len(statuses) != 0 {
		db = db.Where("status IN (?)", statuses)
	}

	if len(productId) != 0 {
		db = db.Where("product_id=?", productId)
	}

	if side != nil {
		db = db.Where("side=?", side)
	}

	if beforeId > 0 {
		db = db.Where("id>?", beforeId)
	}

	if afterId > 0 {
		db = db.Where("id<?", afterId)
	}

	if limit <= 0 {
		limit = 100
	}

	db = db.Order("id DESC").Limit(limit)

	var orders []*models.Order
	err := db.Find(&orders).Error
	return orders, err
}

func (s *Store) AddOrder(order *models.Order) error {
	return s.db.Create(order).Error
}

func (s *Store) UpdateOrder(order *models.Order) error {
	return s.db.Save(order).Error
}

func (s *Store) UpdateOrderStatus(orderId int64, oldStatus, newStatus models.OrderStatus) (bool, error) {
	ret := s.db.Exec("UPDATE g_order SET `status`=? WHERE id=? AND `status`=? ", newStatus, orderId, oldStatus)
	if ret.Error != nil {
		return false, ret.Error
	}
	return ret.RowsAffected > 0, nil
}
