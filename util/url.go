package util

import (
	"fmt"
	"net/url"
	"strconv"
)

func GetProductCandlestickUrl(product_id string, start int64, end int64, granularity string, limit int64) url.URL {
	params := url.Values{}
	params.Add("start", strconv.FormatInt(start, 10))
	params.Add("end", strconv.FormatInt(end, 10))
	params.Add("granularity", granularity)
	params.Add("limit", strconv.FormatInt(limit, 10))

	requestHost := "api.coinbase.com"
	requestPath := fmt.Sprintf("/api/v3/brokerage/products/%s/candles", product_id)

	return url.URL{
		Scheme:   "https",
		Host:     requestHost,
		Path:     requestPath,
		RawQuery: params.Encode(),
	}
}
