package main

import (
	"context"
	"sync"

	"github.com/mmjlee/btc-analysis/api"
	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup

	decimal.MarshalJSONWithoutQuotes = true
	// goroutine to log data from coinbase api to postgres db
	conn := repository.NewCandleConn(ctx)
	defer conn.Conn.Close(ctx)
	candle_logger := client.GetNewAPIClient()
	wg.Add(1)
	go func() {
		defer wg.Done()
		candle_logger.LogCandles(conn)
	}()

	// serve http requests
	candle_pool := repository.NewCandlePool(ctx)
	defer candle_pool.Pool.Close()
	candle_handler := api.NewCandleHandler(candle_pool)
	server := api.GetServer(candle_handler)
	server.ListenAndServe()

	wg.Wait()
}
