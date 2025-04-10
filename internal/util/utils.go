package util

import (
	"fmt"
)

const ISODate = "2006-01-02"

// custom implementation for the pgx CopyFrom
type CandleSliceWithTicker struct {
	Ticker string `json:"ticker"`
	CandleSlice
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
		candle.Start,
		candle.Open,
		candle.High,
		candle.Low,
		candle.Close,
		candle.Volume,
	}, nil
}
