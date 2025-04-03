package util

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var GRANULARITIES = [...]string{"ONE_MINUTE", "FIVE_MINUTE", "FIFTEEN_MINUTE", "THIRTY_MINUTE", "ONE_HOUR", "TWO_HOUR", "SIX_HOUR", "ONE_DAY"}

func ValidateGranularity(granularity string) bool {
	for _, g := range GRANULARITIES {
		if granularity == g {
			return true
		}
	}
	return false
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
