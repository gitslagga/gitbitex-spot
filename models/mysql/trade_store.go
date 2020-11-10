package mysql

import (
	"fmt"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"strings"
)

func (s *Store) GetLastTradeByProductId(productId string) (*models.Trade, error) {
	var trade models.Trade
	err := s.db.Where("product_id =?", productId).Order("id DESC").Limit(1).Find(&trade).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &trade, err
}

func (s *Store) GetTradesByProductId(productId string, count int) ([]*models.Trade, error) {
	db := s.db.Where("product_id =?", productId).Order("id DESC").Limit(count)
	var trades []*models.Trade
	err := db.Find(&trades).Error
	return trades, err
}

func (s *Store) AddTrades(trades []*models.Trade) error {
	if len(trades) == 0 {
		return nil
	}
	var valueStrings []string
	for _, trade := range trades {
		valueString := fmt.Sprintf("('%v', %v, %v, %v, %v, '%v', '%v', %v, %v)",
			trade.ProductId, trade.TakerOrderId, trade.MakerOrderId, trade.Price, trade.Size, trade.Side,
			trade.Time, trade.LogOffset, trade.LogSeq)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT IGNORE  INTO g_trade (product_id,taker_order_id,maker_order_id,"+
		"price,size,side,time,log_offset,log_seq) VALUES %s", strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}
