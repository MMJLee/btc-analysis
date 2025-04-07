package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mmjlee/btc-analysis/api"
	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/mmjlee/btc-analysis/internal/util"
	"github.com/shopspring/decimal"
)

func getServer(c api.CandleHandler) http.Server {
	router := http.NewServeMux()
	router.HandleFunc("GET /candle/{ticker}", c.GetCandles)
	router.HandleFunc("POST /candle/{ticker}", c.GetCandles)
	router.HandleFunc("OPTIONS /candle/{ticker}", c.Options)
	admin_router := http.NewServeMux()
	admin_router.HandleFunc("PUT /candle/{ticker}", c.GetCandles)
	admin_router.HandleFunc("PATCH /candle/{ticker}", c.GetCandles)
	admin_router.HandleFunc("DELETE /candle/{ticker}", c.GetCandles)
	router.Handle("/", util.AuthMiddleware(admin_router))
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))
	middlewares := util.CreateStack(
		util.GzipMiddleware,
		util.CORSMiddleware,
		util.ErrorMiddleware,
		util.LoggingMiddleware,
	)

	return http.Server{
		Addr: "localhost:8080",
		// Addr:    "0.0.0.0:8080",
		Handler: middlewares(v1),
	}
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup

	decimal.MarshalJSONWithoutQuotes = true
	// goroutine to log data from coinbase api to postgres db
	conn := repository.NewCandleConn(ctx)
	defer conn.Conn.Close(ctx)
	candle_logger := client.GetNewAPIClient()
	product_id, limit := "BTC-USD", 3
	wg.Add(1)
	go func() {
		count := 0
		defer wg.Done()
		for {
			now := time.Now().Add(time.Duration(-limit*count) * time.Minute).Truncate(time.Minute)
			start := now.Add(time.Duration(-limit) * time.Minute).Unix()
			end := now.Add(time.Duration(-1) * time.Second).Unix()

			candles_response, err := candle_logger.GetCandles(product_id, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10), strconv.Itoa(limit))
			if err != nil {
				log.Panicf("Error: Main-GetCandles: %v", err)
			}
			// to backfill use limit=350 and time.Sleep(250ms)
			// count++
			// if err := conn.BulkLogCandles(product_id, candles_response.Candles); err != nil {
			if err := conn.InsertCandles(product_id, candles_response.Candles); err != nil {
				log.Panicf("Error: Main-InsertCandles: %v", err)
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()

	// serve http requests
	candle_pool := repository.NewCandlePool(ctx)
	defer candle_pool.Pool.Close()
	candle_handler := api.NewCandleHandler(candle_pool)
	server := getServer(candle_handler)
	server.ListenAndServe()

	wg.Wait()
}
