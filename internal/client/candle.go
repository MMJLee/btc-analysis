package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/database"
)

type CandleResponse struct {
	Candles database.CandleSlice `json:"candles"`
}

func (c CoinbaseClient) getCandles(ticker, start, end, limit string) (CandleResponse, error) {
	var candlesResponse CandleResponse
	candleUrl := getProductCandleUrl(ticker, start, end, "ONE_MINUTE", limit)
	req, err := NewRequest("GET", candleUrl, nil)
	if err != nil {
		return candlesResponse, fmt.Errorf("GetCandles-%w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return candlesResponse, fmt.Errorf("GetCandles-%w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return candlesResponse, fmt.Errorf("GetCandles-%w", err)
	}

	err = json.Unmarshal([]byte(body), &candlesResponse)
	if err != nil {
		return candlesResponse, fmt.Errorf("GetCandles-%w", err)
	}
	return candlesResponse, nil
}

func logRecentCandles(ctx context.Context, client CoinbaseClient, conn database.DBConn, ticker string, limit int) error {
	now := time.Now()
	start := now.Add(time.Duration(-limit)*time.Minute + time.Second).Unix()
	candlesResponse, err := client.getCandles(ticker, strconv.FormatInt(start, 10), strconv.FormatInt(now.Unix(), 10), strconv.Itoa(limit))
	if err != nil {
		return fmt.Errorf("LogRecentCandles-%w", err)
	}
	if err := conn.InsertCandles(ctx, ticker, candlesResponse.Candles); err != nil {
		return fmt.Errorf("LogRecentCandles-%w", err)
	}
	return nil
}

func TrackTicker(ticker string, stopChan chan bool) {
	ctx := context.Background()
	conn := database.NewConn()
	defer conn.Close(ctx)
	client := NewCoinbaseClient()

	t := time.NewTicker(time.Duration(10) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-stopChan:
			log.Println("Stopped tracking", ticker)
		case _ = <-t.C:
			if err := logRecentCandles(ctx, client, conn, ticker, 3); err != nil {
				log.Panic(err)
			}
		}
	}
}

func backfillCandles(ctx context.Context, client CoinbaseClient, conn database.DBConn, ticker string, start, stop, limit int64) error {
	candlesResponse, err := client.getCandles(ticker, strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10), strconv.FormatInt(limit, 10))
	if err != nil {
		return fmt.Errorf("BackfillCandles-%w", err)
	}
	if err := conn.BulkLogCandles(ctx, ticker, candlesResponse.Candles); err != nil {
		return fmt.Errorf("BackfillCandles-%w", err)
	}
	return nil
}

func BackfillTicker(ticker string, start, stop int64, stopChan chan bool) {
	ctx := context.Background()
	conn := database.NewConn()
	defer conn.Close(ctx)
	client := NewCoinbaseClient()
	limit := int64(350)
	for t := start; t < stop; t += limit * 60 {
		select {
		case <-stopChan:
			log.Println("Stopped backfilling", ticker)
		default:
			if err := backfillCandles(ctx, client, conn, ticker, t, (t + (limit * 60) - 1), limit); err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(150) * time.Millisecond)
		}
	}
}
