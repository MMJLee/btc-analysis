package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/repository"
)

type CandleHandler struct {
	repository.CandlePool
}

func NewCandleHandler(repo repository.CandlePool) CandleHandler {
	return CandleHandler{repo}
}

func (c CandleHandler) Options(w http.ResponseWriter, r *http.Request) {
	return
}

func (c CandleHandler) GetCandles(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	if ticker == "" || start == "" || end == "" || limit == "" || offset == "" {
		log.Panic("Error: API-GetCandles: missing required query param")
	}

	candles, err := c.CandlePool.GetCandles(ticker, start, end, limit, offset)
	if err != nil {
		log.Panicf("Error: API-GetCandles-GetCandles: %v", err)
	}
	json_data, err := json.Marshal(candles)
	if err != nil {
		log.Panicf("Error: API-GetCandles-Marshal: %v", err)
	}
	w.Write(json_data)
}

func (c CandleHandler) GetMissingCandles(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	if ticker == "" || start == "" || end == "" || limit == "" || offset == "" {
		log.Panic("Error: API-GetMissingCandles: missing required query param")
	}
	start_int, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	start = strconv.FormatInt(time.Unix(start_int, 0).Truncate(time.Minute).Unix(), 10)

	candles, err := c.CandlePool.GetMissingCandles(ticker, start, end, limit, offset)
	if err != nil {
		log.Panicf("Error: API-GetMissingCandles-GetCandles: %v", err)
	}
	json_data, err := json.Marshal(candles)
	if err != nil {
		log.Panicf("Error: API-GetMissingCandles-Marshal: %v", err)
	}
	w.Write(json_data)
}
