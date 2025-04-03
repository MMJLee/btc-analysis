package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"mjlee.dev/btc-analysis/api"
	"mjlee.dev/btc-analysis/internal/client"
	"mjlee.dev/btc-analysis/internal/util"
)

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	w.Write([]byte("hohoh"))
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup

	conn, err := pgx.Connect(context.Background(), api.DATABASE_CONNECTION_STRING)
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(ctx)

	api_client := client.APIClient{Client: &http.Client{}}
	product_id := "BTC-USD"
	limit := 1
	end := time.Now().Unix()
	start := time.Now().Add(time.Duration(-limit) * time.Minute).Unix()
	granularity := "ONE_MINUTE"

	wg.Add(1)
	go func() {
		defer wg.Done()
		candles_response, err := api_client.GetCandles(product_id, start, end, granularity, limit)
		if err != nil {
			log.Panicf("Error: Main-GetCandles: %v", err)
		}
		_, err = conn.CopyFrom(
			ctx,
			pgx.Identifier{"candle_one_minute"},
			[]string{"ticker", "start", "open", "high", "low", "close", "volume"},
			&util.CandleSliceWithTicker{Ticker: product_id, CandleSlice: candles_response.Candles},
		)
		// if err != nil {
		// 	log.Panicf("Error: Client-LogCandles-CopyFrom: %v", err)
		// }
	}()

	router := http.NewServeMux()
	router.HandleFunc("GET /candle/{product_id}", handle)
	router.HandleFunc("POST /candle/{product_id}", handle)
	router.HandleFunc("OPTIONS /candle/{product_id}", handle)

	admin_router := http.NewServeMux()
	admin_router.HandleFunc("PUT /candle/{product_id}", handle)
	admin_router.HandleFunc("PATCH /candle/{product_id}", handle)
	admin_router.HandleFunc("DELETE /candle/{product_id}", handle)
	router.Handle("/", util.AuthMiddleware(admin_router))

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	middlewares := util.CreateStack(
		util.GzipMiddleware,
		util.CORSMiddleware,
		util.ErrorMiddleware,
		util.LoggingMiddleware,
	)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: middlewares(v1),
	}
	server.ListenAndServe()
	wg.Wait()
}
