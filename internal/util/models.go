package util

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// money represented in cents
const dollar_to_cents uint8 = 100

type StringInt64 int64

func (s StringInt64) Int64() int64 {
	return int64(s)
}

func (s *StringInt64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*s = StringInt64(int64(value))
	return nil
}

type Candle struct {
	Start  StringInt64 `json:"start"`
	Open   string      `json:"open"`
	High   string      `json:"high"`
	Low    string      `json:"low"`
	Close  string      `json:"close"`
	Volume string      `json:"volume"`
}

type CandleWithTicker struct {
	Ticker string      `json:"ticker"`
	Start  StringInt64 `json:"start"`
	Open   string      `json:"open"`
	High   string      `json:"high"`
	Low    string      `json:"low"`
	Close  string      `json:"close"`
	Volume string      `json:"volume"`
}

type CandleSlice []Candle

type CandleResponse struct {
	Candles CandleSlice `json:"candles"`
}

func (c *CandleSlice) Data() any {
	return *c
}

func (c *CandleSlice) Err() error {
	return nil
}

func (c *CandleSlice) Next() bool {
	if len(*c) > 0 {
		return true
	}
	return false
}

func (c *CandleSlice) Values() ([]any, error) {
	if len(*c) > 0 {
		candle := (*c)[0]
		*c = (*c)[1:]
		return []any{
			"BTC-USD",
			candle.Start.Int64(),
			candle.Open,
			candle.High,
			candle.Low,
			candle.Close,
			candle.Volume,
		}, nil
	}
	return nil, fmt.Errorf("Error: CandleSlice-Values")
}
