package mysql

import (
	"fmt"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"strings"
)

func (s *Store) GetLastFillByProductId(productId string) (*models.Fill, error) {
	var fill models.Fill
	err := s.db.Where("product_id =?", productId).Order("id DESC").Limit(1).Find(&fill).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &fill, err
}

func (s *Store) GetUnsettledFillsByOrderId(orderId int64) ([]*models.Fill, error) {
	db := s.db.Where("settled =?", 0).Where("order_id=?", orderId).
		Order("id ASC").Limit(100)

	var fills []*models.Fill
	err := db.Find(&fills).Error
	return fills, err
}

func (s *Store) GetUnsettledFills(count int32) ([]*models.Fill, error) {
	db := s.db.Where("settled =?", 0).Order("id ASC").Limit(count)

	var fills []*models.Fill
	err := db.Find(&fills).Error
	return fills, err
}

func (s *Store) UpdateFill(fill *models.Fill) error {
	return s.db.Save(fill).Error
}

func (s *Store) AddFills(fills []*models.Fill) error {
	if len(fills) == 0 {
		return nil
	}
	var valueStrings []string
	for _, fill := range fills {
		valueString := fmt.Sprintf("('%v', %v, %v, %v, %v,%v, %v,'%v',%v,%v,'%v',%v,'%v',%v,%v)",
			fill.ProductId, fill.TradeId, fill.OrderId, fill.MessageSeq, fill.Size, fill.Price, fill.Funds,
			fill.Liquidity, fill.Fee, fill.Settled, fill.Side, fill.Done, fill.DoneReason, fill.LogOffset, fill.LogSeq)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT IGNORE INTO g_fill (product_id,trade_id,order_id, message_seq,size,"+
		"price,funds,liquidity,fee,settled,side,done,done_reason,log_offset,log_seq) VALUES %s",
		strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}
