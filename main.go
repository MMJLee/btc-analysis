package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

type RecipesHandler struct{}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my recipe page"))
}

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(context.Background(), api.DATABASE_CONNECTION_STRING)
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// data.getCandlesticks(conn, ctx)

	mux := http.NewServeMux()
	mux.Handle("/recipes/", &RecipesHandler{})
	go http.ListenAndServe("localhost:8080", mux)

	client := client.APIClient{Client: &http.Client{}}

	start := time.Now().Add(time.Duration(-24) * time.Hour).Unix()
	fmt.Println(start)
	end := time.Unix(start, 0).Add(time.Duration(350) * time.Minute).Unix()
	fmt.Println(end)
	// go client.logCandlesticks(("BTC-USD", start, end, "ONE_MINUTE", 10))
	candlesticks, err := client.GetCandlesticks("BTC-USD", start, end, "ONE_MINUTE", 350)
	if err != nil {
		log.Fatal(err)
	}
	for _, candlestick := range candlesticks {
		fmt.Printf("%+v\n", candlestick)
	}
}
