package main

import (
	"flag"
	"sync"

	"github.com/mmjlee/btc-analysis/api"
	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/database"
	"github.com/shopspring/decimal"
)

func main() {
	tickerPtr := flag.String("ticker", "BTC-USD", "Trading pair")
	modePtr := flag.String("mode", "log", "Startup mode - log/fill/both")
	startPtr := flag.Int64("start", 0, "Start unix epoch if fill/both")
	endPtr := flag.Int64("end", 0, "End unix epoch if fill/both")
	flag.Parse()

	var wg sync.WaitGroup
	decimal.MarshalJSONWithoutQuotes = true
	trackMap := make(map[string]chan bool)
	trackMut := new(sync.Mutex)
	backfillMap := make(map[string]chan bool)
	backfillMut := new(sync.Mutex)

	ticker := *tickerPtr
	mode := *modePtr
	start := *startPtr
	end := *endPtr

	switch mode {
	case "log":
		wg.Add(1)
		stopChan := make(chan bool)
		trackMap[ticker] = stopChan
		go func() {
			defer wg.Done()
			client.TrackTicker(ticker, stopChan)
		}()
	case "fill":
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.BackfillTicker(ticker, start, end, nil)
			backfillMut.Lock()
			delete(backfillMap, ticker)
			backfillMut.Unlock()
		}()
	case "both":
		wg.Add(2)
		stopChan := make(chan bool)
		trackMap[ticker] = stopChan
		go func() {
			defer wg.Done()
			client.TrackTicker(ticker, stopChan)
		}()
		go func() {
			defer wg.Done()
			client.BackfillTicker(ticker, start, end, nil)
			backfillMut.Lock()
			delete(backfillMap, ticker)
			backfillMut.Unlock()
		}()
	}

	// serve http requests
	dbPool := database.NewPool()
	defer dbPool.Close()
	redis := database.NewRedis()
	defer redis.Close()

	candleHandler := api.NewCandleHandler(dbPool)
	trackHandler := api.NewTrackHandler(dbPool, trackMap, trackMut)
	backfillHandler := api.NewBackfillHandler(dbPool, backfillMap, backfillMut)
	authHandler := api.NewAuthHandler(dbPool, redis)
	handlers := []api.Handler{candleHandler, trackHandler, backfillHandler, authHandler}
	server := api.GetServer(redis, handlers...)
	server.ListenAndServe()

	wg.Wait()
}
