package mysql

import (
	"fmt"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/jinzhu/gorm"
	"strings"
)

func (s *Store) GetTicksByProductId(productId string, granularity int64, limit int) ([]*models.Tick, error) {
	db := s.db.Where("product_id =?", productId).Where("granularity=?", granularity).
		Order("time DESC").Limit(limit)
	var ticks []*models.Tick
	err := db.Find(&ticks).Error
	return ticks, err
}

func (s *Store) GetLastTickByProductId(productId string, granularity int64) (*models.Tick, error) {
	var tick models.Tick
	err := s.db.Raw("SELECT * FROM g_tick WHERE product_id=? AND granularity=? ORDER BY time DESC LIMIT 1",
		productId, granularity).Scan(&tick).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &tick, err
}

func (s *Store) AddTicks(ticks []*models.Tick) error {
	if len(ticks) == 0 {
		return nil
	}
	var valueStrings []string
	for _, tick := range ticks {
		valueString := fmt.Sprintf("('%v', %v, %v, %v, %v, %v, %v, %v,%v,%v)",
			tick.ProductId, tick.Granularity, tick.Time, tick.Open, tick.Low, tick.High, tick.Close,
			tick.Volume, tick.LogOffset, tick.LogSeq)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("REPLACE INTO g_tick (product_id,granularity,time,open,low,high,close,"+
		"volume,log_offset,log_seq) VALUES %s", strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}
