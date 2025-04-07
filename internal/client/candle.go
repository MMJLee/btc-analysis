package client

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/mmjlee/btc-analysis/internal/util"
)

func (a *APIClient) GetCandles(product_id, start, end, limit string) (util.CandleResponse, error) {
	candle_url := util.GetProductCandleUrl(product_id, start, end, "ONE_MINUTE", limit)
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
		log.Panicf("Error: Client-GetCandles-Response: %v", string(body))
	}

	var candles_response util.CandleResponse

	err = json.Unmarshal([]byte(body), &candles_response)
	if err != nil {
		log.Panicf("Error: Client-GetCandles-Unmarshal: %v", err)
	}
	return candles_response, nil
}

func (a *APIClient) LogCandles(conn repository.CandleConn) (util.CandleResponse, error) {
	product_id, limit, count := "BTC-USD", 3, 0
	for {
		now := time.Now().Add(time.Duration(-limit*count) * time.Minute).Truncate(time.Minute)
		start := now.Add(time.Duration(-limit) * time.Minute).Unix()
		end := now.Add(time.Duration(-1) * time.Second).Unix()

		candles_response, err := a.GetCandles(product_id, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10), strconv.Itoa(limit))
		if err != nil {
			log.Panicf("Error: Main-GetCandles: %v", err)
		}
		// to backfill, use limit=350, sleep=250ms, count++, and BulkLogCandles
		// if err := conn.BulkLogCandles(product_id, candles_response.Candles); err != nil {
		if err := conn.InsertCandles(product_id, candles_response.Candles); err != nil {
			log.Panicf("Error: Main-InsertCandles: %v", err)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}
