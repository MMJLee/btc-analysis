package main

import (
	"sync"

	"github.com/mmjlee/btc-analysis/api"
	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/shopspring/decimal"
)

func main() {
	var wg sync.WaitGroup
	var mut *sync.Mutex
	// ctxt, cnclFn := context.WithCancel()
	decimal.MarshalJSONWithoutQuotes = true

	// goroutine to log data from coinbase api to postgres db
	// TrackTicker can also be called from the api
	tickerMap := make(map[string]chan bool)
	ticker := "BTC-USD"
	stopChan := make(chan bool)
	tickerMap[ticker] = stopChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.TrackTicker(ticker, stopChan)
	}()

	// serve http requests
	dbPool := repository.NewPool()
	defer dbPool.Pool.Close()
	candleHandler := api.NewCandleHandler(dbPool)
	trackHandler := api.NewTrackHandler(dbPool, tickerMap, mut)
	server := api.GetServer(candleHandler, trackHandler)
	server.ListenAndServe()

	wg.Wait()
}
