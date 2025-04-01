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

type CandleSlice []Candle

type CandleSliceWithTicker struct {
	Ticker string `json:"ticker"`
	CandleSlice
}
type CandleResponse struct {
	Candles CandleSlice `json:"candles"`
}

func (c *CandleSliceWithTicker) Data() any {
	return *c
}

func (c *CandleSliceWithTicker) Err() error {
	return nil
}

func (c *CandleSliceWithTicker) Next() bool {
	if len(c.CandleSlice) == 0 {
		return false
	}
	return true
}

func (c *CandleSliceWithTicker) Values() ([]any, error) {
	if len(c.CandleSlice) == 0 {
		return nil, fmt.Errorf("Error: CandleSlice-Values")
	}
	candle := (c.CandleSlice)[0]
	*&c.CandleSlice = (c.CandleSlice[1:])
	return []any{
		c.Ticker,
		candle.Start.Int64(),
		candle.Open,
		candle.High,
		candle.Low,
		candle.Close,
		candle.Volume,
	}, nil
}
