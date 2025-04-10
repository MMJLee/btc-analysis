package main

import (
	"sync"

	"github.com/mmjlee/btc-analysis/api"
	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/database"
	"github.com/shopspring/decimal"
)

func main() {
	var wg sync.WaitGroup
	decimal.MarshalJSONWithoutQuotes = true

	ticker := "BTC-USD"
	tickerMap := make(map[string]chan bool)
	stopChan := make(chan bool)
	tickerMap[ticker] = stopChan
	wg.Add(1)
	go func() {
		defer wg.Done()
		// TrackTicker can also be started with the api as a goroutine /track/{ticker}
		client.TrackTicker(ticker, stopChan)
	}()

	// serve http requests
	dbPool := database.NewPool()
	defer dbPool.Close()
	candleHandler := api.NewCandleHandler(dbPool)
	trackMut := new(sync.Mutex)
	trackHandler := api.NewTrackHandler(dbPool, tickerMap, trackMut)
	server := api.GetServer(candleHandler, trackHandler)
	server.ListenAndServe()

	wg.Wait()
}
