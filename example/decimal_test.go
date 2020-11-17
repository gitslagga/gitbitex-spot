package example

import (
	"github.com/shopspring/decimal"
	"github.com/siddontang/go-log/log"
	"testing"
)

func TestDecimalToFloat64(t *testing.T) {
	num1, err := decimal.NewFromString("123456")
	if err != nil {
		log.Error(num1, err)
		return
	}

	num2, exact := num1.Float64()
	log.Infoln(num2, exact == true) //123456 true

	num3, err := decimal.NewFromString("123456.12")
	if err != nil {
		log.Error(num3, err)
		return
	}

	num4, exact := num3.Float64()
	log.Infoln(num4, exact == false) //123456.12 true
}

func TestDecimalIsZero(t *testing.T) {
	num1 := decimal.Decimal{}

	log.Infoln(num1, num1.Equal(decimal.New(0, 0))) //0 true
	log.Infoln(num1, num1.Equal(decimal.New(0, 1))) //0 true
	log.Infoln(decimal.New(1, 0))                   //1
}
