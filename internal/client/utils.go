package client

import (
	"fmt"
	"net/url"
)

func getProductCandleUrl(ticker, start, end, granularity, limit string) url.URL {
	params := url.Values{}
	params.Add("start", start)
	params.Add("end", end)
	params.Add("granularity", granularity)
	params.Add("limit", limit)

	requestHost := "api.coinbase.com"
	requestPath := fmt.Sprintf("/api/v3/brokerage/products/%s/candles", ticker)

	return url.URL{
		Scheme:   "https",
		Host:     requestHost,
		Path:     requestPath,
		RawQuery: params.Encode(),
	}
}
