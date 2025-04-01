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
)

// 	now := time.Now().Unix()
// 	if start > now {
// 		return nil, fmt.Errorf("Error: candlestick time is in the future data")
// 	}

// granularities := []string{"ONE_MINUTE", "FIVE_MINUTE", "FIFTEEN_MINUTE", "THIRTY_MINUTE", "ONE_HOUR", "TWO_HOUR", "SIX_HOUR", "ONE_DAY"}
//
//	if !slices.Contains(granularities, granularity) {
//		granularity = "UNKNOWN_GRANULARITY"
//	}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("reasdfasdf")
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
	// data.getCandlesticks(conn, ctx)

	// router := http.NewServeMux()
	// router.HandleFunc("GET /candle/{product_id}", handle)
	// router.HandleFunc("POST /candle/{product_id}", handle)
	// router.HandleFunc("PUT /candle/{product_id}", handle)
	// router.HandleFunc("PATCH /candle/{product_id}", handle)
	// router.HandleFunc("DELETE /candle/{product_id}", handle)
	// router.HandleFunc("OPTIONS /candle/{product_id}", handle)

	// server := http.Server{
	// 	Addr:    "localhost:8080",
	// 	Handler: Logging(router),
	// }
	// server.ListenAndServe()

	api_client := client.APIClient{Client: &http.Client{}}

	limit := 2
	end := time.Now().Unix()
	start := time.Now().Add(time.Duration(-limit) * time.Minute).Unix()

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.LogCandles(ctx, conn, &api_client, "BTC-USD", start, end, "ONE_MINUTE", limit)
	}()
	wg.Wait()

}
