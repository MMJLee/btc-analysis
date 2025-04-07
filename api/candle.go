package api

import (
	"encoding/json"
	"log"
	"net/http"

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
		log.Panic("Error: API-GetProduct: missing required query param")
	}

	candles, err := c.CandlePool.GetCandles(ticker, start, end, limit, offset)
	if err != nil {
		log.Panicf("Error: API-GetProduct-GetCandles: %v", err)
	}
	json_data, err := json.Marshal(candles)
	if err != nil {
		log.Panicf("Error: API-GetProduct-Marshal: %v", err)
	}
	w.Write(json_data)
}
