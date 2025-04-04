package util

import (
	"fmt"
	"net/url"
)

func GetProductCandleUrl(product_id, start, end, granularity, limit string) url.URL {
	params := url.Values{}
	params.Add("start", start)
	params.Add("end", start)
	params.Add("granularity", granularity)
	params.Add("limit", limit)

	requestHost := "api.coinbase.com"
	requestPath := fmt.Sprintf("/api/v3/brokerage/products/%s/candles", product_id)

	return url.URL{
		Scheme:   "https",
		Host:     requestHost,
		Path:     requestPath,
		RawQuery: params.Encode(),
	}
}
