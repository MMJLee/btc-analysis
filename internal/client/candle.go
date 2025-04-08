package client

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/mmjlee/btc-analysis/internal/util"
)

func (a *APIClient) GetCandles(product_id, start, end, limit string) (util.CandleResponse, error) {
	var candles_response util.CandleResponse
	candle_url := util.GetProductCandleUrl(product_id, start, end, "ONE_MINUTE", limit)
	req, err := a.NewRequest("GET", candle_url, nil)
	if err != nil {
		return candles_response, util.WrappedError{Err: err, Message: "Client-GetCandles-NewRequest"}
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return candles_response, util.WrappedError{Err: err, Message: "Client-GetCandles-Do"}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return candles_response, util.WrappedError{Err: err, Message: "Client-GetCandles-ReadAll"}
	}

	if resp.StatusCode != 200 {
		return candles_response, util.WrappedError{Err: errors.New("something went wrong"), Message: "Client-GetCandles-Response"}
	}

	err = json.Unmarshal([]byte(body), &candles_response)
	if err != nil {
		return candles_response, util.WrappedError{Err: err, Message: "Client-GetCandles-Unmarshal"}
	}
	return candles_response, nil
}

func (a *APIClient) LogCandles(conn repository.CandleConn, product_id string) {
	limit, count := 3, 0
	for {
		now := time.Now().Truncate(time.Minute)
		start := now.Add(time.Duration(-limit) * time.Minute).Unix()
		end := now.Add(time.Duration(-1) * time.Second).Unix()
		count++

		candles_response, err := a.GetCandles(product_id, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10), strconv.Itoa(limit))
		if err != nil {
			log.Panic(util.WrappedError{Err: err, Message: "Client-LogCandles-GetCandles"}.Error())
		}
		if err := conn.InsertCandles(product_id, candles_response.Candles); err != nil {
			log.Panic(util.WrappedError{Err: err, Message: "Client-LogCandles-InsertCandles"}.Error())
		}
		if count > 5 {
			count = 0
			log.Println("Logging:", product_id, start)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func (a *APIClient) BackfillCandles(conn repository.CandleConn, product_id string, start, stop int64) {
	limit, count := 350, 0
	now := time.Unix(start, 0)
	for {
		start := now.Add(time.Duration(count*limit) * time.Minute).Unix()
		end := now.Add(time.Duration((count+1)*limit)*time.Minute - time.Second).Unix()
		count++
		candles_response, err := a.GetCandles(product_id, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10), strconv.Itoa(limit))
		if err != nil {
			log.Panic(util.WrappedError{Err: err, Message: "Client-BackfillCandles-GetCandles"}.Error())
		}
		if err := conn.BulkLogCandles(product_id, candles_response.Candles); err != nil {
			log.Panic(util.WrappedError{Err: err, Message: "Client-BackfillCandles-BulkLogCandles"}.Error())
		}
		time.Sleep(time.Duration(150) * time.Millisecond)
		if start > stop {
			return
		}
	}
}
