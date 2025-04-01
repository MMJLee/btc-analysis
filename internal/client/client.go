package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/jackc/pgx/v5"
	"mjlee.dev/btc-analysis/api"
	"mjlee.dev/btc-analysis/internal/util"
)

type APIClient struct {
	*http.Client
}

func (a APIClient) NewRequest(method string, url url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		log.Fatalf("Error: Client-NewRequest-NewRequest: %v", err)
	}

	jwt, err := api.BuildJWT(method, url.Host, url.Path)
	if err != nil {
		log.Fatalf("Error: Client-NewRequest-BuildJWT: %v", err)
	}

	bearer := "Bearer " + jwt
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (a *APIClient) GetCandles(product_id string, start int64, end int64, granularity string, limit int) (util.CandleResponse, error) {
	candle_url := util.GetProductCandleUrl(product_id, start, end, granularity, limit)
	req, err := a.NewRequest("GET", candle_url, nil)
	if err != nil {
		log.Fatalf("Error: Client-GetCandles-NewRequest: %v", err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Fatalf("Error: Client-GetCandles-Do: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: Client-GetCandles-ReadAll: %v", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Request succeeded!")
	} else {
		fmt.Println("Request failed with status:", resp.StatusCode)
		panic(body)
	}

	var candles_response util.CandleResponse

	err = json.Unmarshal([]byte(body), &candles_response)
	if err != nil {
		log.Fatalf("Error: GetCandles-Unmarshal: %v", err)
	}

	return candles_response, nil
}

func LogCandles(ctx context.Context, conn *pgx.Conn, a *APIClient, product_id string, start int64, end int64, granularity string, limit int) {
	candles_response, err := a.GetCandles(product_id, start, end, granularity, limit)
	if err != nil {
		log.Fatalf("Error: Client-LogCandles-GetCandles: %v", err)
	}

	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"candle_one_minute"},
		[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
		&util.CandleSliceWithTicker{Ticker: product_id, CandleSlice: candles_response.Candles},
	)
	if err != nil {
		log.Fatalf("Error: Client-LogCandles-CopyFrom: %v", err)
	}
}
