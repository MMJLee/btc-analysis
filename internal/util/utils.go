package util

import (
	"fmt"

	decimal "github.com/jackc/pgx-shopspring-decimal"
)

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
		int64(candle.Start),
		decimal.Decimal(candle.Open),
		decimal.Decimal(candle.High),
		decimal.Decimal(candle.Low),
		decimal.Decimal(candle.Close),
		float64(candle.Volume),
	}, nil
}
