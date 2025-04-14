package database

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const ISODate = "2006-01-02"

type StringInt64 int64

func (s *StringInt64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("UnmarshalJSON-%w", err)
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("UnmarshalJSON-%w", err)
	}
	*s = StringInt64(value)
	return nil
}

type StringFloat64 float64

func (s *StringFloat64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("UnmarshalJSON-%w", err)
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("UnmarshalJSON-%w", err)
	}
	*s = StringFloat64(value)
	return nil
}

// custom implementation for the pgx CopyFrom
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
		return nil, fmt.Errorf("CandleSliceWithTicker-Values")
	}
	candle := (c.CandleSlice)[0]
	*&c.CandleSlice = (c.CandleSlice[1:])
	return []any{
		c.Ticker,
		candle.Start,
		candle.Open,
		candle.High,
		candle.Low,
		candle.Close,
		candle.Volume,
	}, nil
}
