package util

import (
	"encoding/json"
	"strconv"

	"github.com/shopspring/decimal"
)

type StringInt64 int64

func (s *StringInt64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*s = StringInt64(value)
	return nil
}

type StringDecimal decimal.Decimal

func (s *StringDecimal) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := decimal.NewFromString(str)
	if err != nil {
		return err
	}
	*s = StringDecimal(value)
	return nil
}

type StringFloat64 float64

func (s *StringFloat64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*s = StringFloat64(value)
	return nil
}

type Candle struct {
	Ticker string        `json:"ticker"`
	Start  StringInt64   `json:"start"`
	Open   StringDecimal `json:"open"`
	High   StringDecimal `json:"high"`
	Low    StringDecimal `json:"low"`
	Close  StringDecimal `json:"close"`
	Volume StringFloat64 `json:"volume"`
}

type CandleSlice []Candle

type CandleSliceWithTicker struct {
	Ticker string `json:"ticker"`
	CandleSlice
}
type CandleResponse struct {
	Candles CandleSlice `json:"candles"`
}
