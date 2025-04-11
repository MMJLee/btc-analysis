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
	tickerMap := make(map[string]chan bool)
	ticker := *tickerPtr

	switch *modePtr {
	case "log":
		wg.Add(1)
		stopChan := make(chan bool)
		tickerMap[ticker] = stopChan
		go func() {
			defer wg.Done()
			client.TrackTicker(ticker, stopChan)
		}()
	case "fill":
		wg.Add(1)
		go func() {
			defer wg.Done()
			go client.BackfillTicker(ticker, *startPtr, *endPtr, nil)
		}()
	case "both":
		wg.Add(2)
		stopChan := make(chan bool)
		tickerMap[ticker] = stopChan
		go func() {
			defer wg.Done()
			client.TrackTicker(ticker, stopChan)
		}()
		go func() {
			defer wg.Done()
			go client.BackfillTicker(ticker, *startPtr, *endPtr, nil)
		}()
	}

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
