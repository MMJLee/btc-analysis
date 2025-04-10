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
		return WrappedError{Err: err, Message: "Util-StringInt64-UnmarshalJSON-Unmarshal"}
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return WrappedError{Err: err, Message: "Util-StringInt64-UnmarshalJSON-ParseInt"}
	}
	*s = StringInt64(value)
	return nil
}

type StringFloat64 float64

func (s *StringFloat64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return WrappedError{Err: err, Message: "Util-StringFloat64-UnmarshalJSON-Unmarshal"}
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return WrappedError{Err: err, Message: "Util-StringFloat64-UnmarshalJSON-ParseInt"}
	}
	*s = StringFloat64(value)
	return nil
}

type Candle struct {
	Ticker string          `json:"ticker"`
	Start  StringInt64     `json:"start"`
	Open   decimal.Decimal `json:"open"`
	High   decimal.Decimal `json:"high"`
	Low    decimal.Decimal `json:"low"`
	Close  decimal.Decimal `json:"close"`
	Volume StringFloat64   `json:"volume"`
}

type CandleSlice []Candle

type CandleResponse struct {
	Candles CandleSlice `json:"candles"`
}
