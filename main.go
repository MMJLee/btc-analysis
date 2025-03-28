package main

import (
	"fmt"
	"log"
	"net/http"

	"mjlee.dev/btc-analysis/api"
)

// 	now := time.Now().Unix()
// 	if start > now {
// 		return nil, fmt.Errorf("Error: candlestick time is in the future data")
// 	}

// granularities := []string{"ONE_MINUTE", "FIVE_MINUTE", "FIFTEEN_MINUTE", "THIRTY_MINUTE", "ONE_HOUR", "TWO_HOUR", "SIX_HOUR", "ONE_DAY"}
//
//	if !slices.Contains(granularities, granularity) {
//		granularity = "UNKNOWN_GRANULARITY"
//	}

func main() {
	client := api.APIClient{Client: &http.Client{}}
	candlesticks, err := client.GetCandlesticks("BTC-USD", 1735711200, 1735711320, "ONE_MINUTE", 2)
	if err != nil {
		log.Fatal(err)
	}
	for _, candlestick := range candlesticks {
		fmt.Printf("%+v\n", candlestick)
	}
}
