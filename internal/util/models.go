package util

type StringInt64 int64

type Candle struct {
	Ticker string `json:"ticker"`
	Start  string `json:"start"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

type CandleSlice []Candle

type CandleSliceWithTicker struct {
	Ticker string `json:"ticker"`
	CandleSlice
}
type CandleResponse struct {
	Candles CandleSlice `json:"candles"`
}
