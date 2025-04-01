package util

import (
	"fmt"
	"net/url"
	"strconv"
)

func GetProductCandleUrl(product_id string, start int64, end int64, granularity string, limit int) url.URL {
	params := url.Values{}
	params.Add("start", strconv.FormatInt(start, 10))
	params.Add("end", strconv.FormatInt(end, 10))
	params.Add("granularity", granularity)
	params.Add("limit", strconv.Itoa(limit))

	requestHost := "api.coinbase.com"
	requestPath := fmt.Sprintf("/api/v3/brokerage/products/%s/candles", product_id)

	return url.URL{
		Scheme:   "https",
		Host:     requestHost,
		Path:     requestPath,
		RawQuery: params.Encode(),
	}
}
