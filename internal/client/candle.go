package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/mmjlee/btc-analysis/internal/util"
)

func (c CoinbaseClient) GetCandles(ticker, start, end, limit string) (util.CandleResponse, error) {
	var candlesResponse util.CandleResponse
	candleUrl := util.GetProductCandleUrl(ticker, start, end, "ONE_MINUTE", limit)
	req, err := NewRequest("GET", candleUrl, nil)
	if err != nil {
		return candlesResponse, util.WrappedError{Err: err, Message: "Client-GetCandles-NewRequest"}
	}
	resp, err := c.Do(req)
	if err != nil {
		return candlesResponse, util.WrappedError{Err: err, Message: "Client-GetCandles-Do"}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return candlesResponse, util.WrappedError{Err: err, Message: "Client-GetCandles-ReadAll"}
	}

	err = json.Unmarshal([]byte(body), &candlesResponse)
	if err != nil {
		return candlesResponse, util.WrappedError{Err: err, Message: "Client-GetCandles-Unmarshal"}
	}
	return candlesResponse, nil
}

func LogRecentCandles(ctx context.Context, client CoinbaseClient, conn repository.DBConn, ticker string, limit int) error {
	now := time.Now()
	start := now.Add(time.Duration(-limit)*time.Minute + time.Second).Unix()
	candlesResponse, err := client.GetCandles(ticker, strconv.FormatInt(start, 10), strconv.FormatInt(now.Unix(), 10), strconv.Itoa(limit))
	if err != nil {
		return util.WrappedError{Err: err, Message: "Client-LogCandles-GetCandles"}
	}
	if err := conn.InsertCandles(ctx, ticker, candlesResponse.Candles); err != nil {
		return util.WrappedError{Err: err, Message: "Client-LogCandles-InsertCandles"}
	}
	return nil
}

func TrackTicker(ticker string, stopChan chan bool) error {
	ctx := context.Background()
	conn := repository.NewConn()
	defer conn.Close(ctx)
	client := NewCoinbaseClient()
	limit := 3

	for {
		select {
		case <-stopChan:
			log.Println("Stopped tracking", ticker)
			return nil
		default:
			if err := LogRecentCandles(ctx, client, conn, ticker, limit); err != nil {
				return util.WrappedError{Err: err, Message: "Client-TrackTicker"}
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}
}

func BackfillCandles(ctx context.Context, client CoinbaseClient, conn repository.DBConn, ticker string, start, stop, limit int64) error {
	candlesResponse, err := client.GetCandles(ticker, strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10), strconv.FormatInt(limit, 10))
	if err != nil {
		return util.WrappedError{Err: err, Message: "Client-BackfillCandles-GetCandles"}
	}
	if err := conn.BulkLogCandles(ctx, ticker, candlesResponse.Candles); err != nil {
		return util.WrappedError{Err: err, Message: "Client-BackfillCandles-BulkLogCandles"}
	}
	return nil
}

func BackfillTicker(ticker string, start, stop int64, stopChan chan bool) error {
	ctx := context.Background()
	conn := repository.NewConn()
	defer conn.Close(ctx)
	client := NewCoinbaseClient()
	limit := int64(350)
	for t := start; t < stop; t += limit * 60 {
		select {
		case <-stopChan:
			log.Println("Stopped backfilling", ticker)
			return nil
		default:
			if err := BackfillCandles(ctx, client, conn, ticker, t, (t + (limit * 60) - 1), limit); err != nil {
				return util.WrappedError{Err: err, Message: "Client-BackfillTicker-BackfillCandles"}
			}
			time.Sleep(time.Duration(150) * time.Millisecond)
		}
	}
	return nil
}
