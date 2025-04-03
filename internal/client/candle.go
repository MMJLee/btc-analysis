package client

import (
	"encoding/json"
	"io"
	"log"

	"mjlee.dev/btc-analysis/internal/util"
)

func (a *APIClient) GetCandles(product_id string, start int64, end int64, granularity string, limit int) (util.CandleResponse, error) {
	candle_url := util.GetProductCandleUrl(product_id, start, end, granularity, limit)
	req, err := a.NewRequest("GET", candle_url, nil)
	if err != nil {
		log.Panicf("Error: Client-GetCandles-NewRequest: %v", err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Panicf("Error: Client-GetCandles-Do: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("Error: Client-GetCandles-ReadAll: %v", err)
	}

	if resp.StatusCode == 200 {
		log.Printf("Info: Client-GetCandles-Response %v:%v", product_id, start)
	} else {
		log.Panicf("Error: Client-GetCandles-Response: %v", body)
	}

	var candles_response util.CandleResponse

	err = json.Unmarshal([]byte(body), &candles_response)
	if err != nil {
		log.Panicf("Error: Client-GetCandles-Unmarshal: %v", err)
	}

	return candles_response, nil
}
