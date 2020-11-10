package mysql

import (
	"fmt"
	"github.com/gitslagga/gitbitex-spot/models"
	"strings"
)

func (s *Store) GetUnsettledBillsByUserId(userId int64, currency string) ([]*models.Bill, error) {
	db := s.db.Where("settled =?", 0).Where("user_id=?", userId).
		Where("currency=?", currency).Order("id ASC").Limit(100)

	var bills []*models.Bill
	err := db.Find(&bills).Error
	return bills, err
}

func (s *Store) GetUnsettledBills() ([]*models.Bill, error) {
	db := s.db.Where("settled =?", 0).Order("id ASC").Limit(100)

	var bills []*models.Bill
	err := db.Find(&bills).Error
	return bills, err
}

func (s *Store) AddBills(bills []*models.Bill) error {
	if len(bills) == 0 {
		return nil
	}
	var valueStrings []string
	for _, bill := range bills {
		valueString := fmt.Sprintf("(%v, '%v', %v, %v, '%v', %v, '%v')",
			bill.UserId, bill.Currency, bill.Available, bill.Hold, bill.Type, bill.Settled, bill.Notes)
		valueStrings = append(valueStrings, valueString)
	}
	sql := fmt.Sprintf("INSERT INTO g_bill (user_id,currency,available,hold,type,settled,notes) VALUES %s", strings.Join(valueStrings, ","))
	return s.db.Exec(sql).Error
}

func (s *Store) UpdateBill(bill *models.Bill) error {
	return s.db.Save(bill).Error
}
