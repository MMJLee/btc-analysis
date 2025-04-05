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
)

func test(ctx context.Context) {
	candle_pool, err := repository.NewCandlePool(ctx)
	defer candle_pool.Pool.Close()

	if err != nil {
		log.Panicf("Error: Main: %v", err)
	}
	candle_handler := api.NewCandleHandler(candle_pool)

	router := http.NewServeMux()
	router.HandleFunc("GET /candle/{product_id}", candle_handler.GetProduct)
	router.HandleFunc("POST /candle/{product_id}", candle_handler.GetProduct)
	router.HandleFunc("OPTIONS /candle/{product_id}", candle_handler.GetProduct)

	admin_router := http.NewServeMux()
	admin_router.HandleFunc("PUT /candle/{product_id}", candle_handler.GetProduct)
	admin_router.HandleFunc("PATCH /candle/{product_id}", candle_handler.GetProduct)
	admin_router.HandleFunc("DELETE /candle/{product_id}", candle_handler.GetProduct)
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
}

func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		test(ctx)
	}()

	// conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_CONNECTION_STRING"))

	conn, err := repository.NewCandleConn(ctx)
	defer conn.Conn.Close(ctx)

	if err != nil {
		log.Panicf("Error: Main: %v", err)
	}
	candle_logger := client.GetNewAPIClient()

	product_id := "BTC-USD"
	limit := 1

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			now := time.Now().Truncate(time.Minute)
			start := now.Add(time.Duration(-limit) * time.Minute).Unix()
			end := now.Add(time.Duration(-limit) * time.Second).Unix()

			candles_response, err := candle_logger.GetCandles(product_id, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10), strconv.Itoa(limit))
			if err != nil {
				log.Panicf("Error: Main-GetCandles: %v", err)
			}

			if err := conn.CopyCandles(product_id, candles_response.Candles); err != nil {
				log.Panicf("Error: Client-LogCandles-CopyFrom: %v", err)
			}
			time.Sleep(time.Duration(70) * time.Second)
		}
	}()

	wg.Wait()
}
