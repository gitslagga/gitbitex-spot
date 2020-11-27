package utils

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

func DecimalAscComparator(a, b interface{}) int {
	aAsserted := a.(decimal.Decimal)
	bAsserted := b.(decimal.Decimal)
	return aAsserted.Cmp(bAsserted)
}

func DecimalDescComparator(a, b interface{}) int {
	aAsserted := a.(decimal.Decimal)
	bAsserted := b.(decimal.Decimal)
	return bAsserted.Cmp(aAsserted)
}

func StartPosOfTime(unixTime int64, granularity int64) int64 {
	return unixTime / (granularity * 60) * (granularity * 60)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func AToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func F64ToA(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func I64ToA(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func IntToA(i int) string {
	return strconv.Itoa(i)
}

func DToF64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func MinInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func SnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return string(out)
}

// struct to map
func StructToMapViaJson(data interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &m)
	return m
}

// map to struct
func MapToStructViaJson(data map[string]interface{}) interface{} {
	var m interface{}
	j, _ := json.Marshal(data)
	_ = json.Unmarshal(j, &m)
	return m
}

// get orderSN
func GetOrderSN() string {
	letterByte := []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`)
	n := 6

	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterByte[r.Intn(len(letterByte))]
	}

	return fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), b)
}
