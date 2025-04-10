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
	trackMut := new(sync.Mutex)
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
	dbPool := database.NewPool()
	defer dbPool.Close()
	candleHandler := api.NewCandleHandler(dbPool)
	trackHandler := api.NewTrackHandler(dbPool, tickerMap, trackMut)
	server := api.GetServer(candleHandler, trackHandler)
	server.ListenAndServe()

	wg.Wait()
}
