package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/repository"
	"github.com/mmjlee/btc-analysis/internal/util"
)

type CandleHandler struct {
	repository.DBPool
}

func NewCandleHandler(repo repository.DBPool) CandleHandler {
	return CandleHandler{repo}
}

func (c CandleHandler) Get(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	queryParams := r.URL.Query()
	start := queryParams.Get("start")
	end := queryParams.Get("end")
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")
	if ticker == "" || start == "" || end == "" || limit == "" || offset == "" {
		log.Panic("Error: API-GetCandles: missing required query param")
	}
	missing, _ := strconv.ParseBool(queryParams.Get("missing"))
	if missing { //make sure start epoch is truncated to the minute
		startInt, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		start = strconv.FormatInt(time.Unix(startInt, 0).Truncate(time.Minute).Unix(), 10)
	}
	candles, err := c.DBPool.GetCandles(r.Context(), ticker, start, end, limit, offset, missing)
	if err != nil {
		log.Panicf("Error: API-GetCandles-GetCandles: %v", err)
	}
	jsonData, err := json.Marshal(candles)
	if err != nil {
		log.Panicf("Error: API-GetCandles-Marshal: %v", err)
	}
	w.Write(jsonData)
}

func (c CandleHandler) Post(w http.ResponseWriter, r *http.Request) {
	util.WriteError(w, http.StatusNotImplemented)
	return
}

func (c CandleHandler) Put(w http.ResponseWriter, r *http.Request) {
	util.WriteError(w, http.StatusNotImplemented)
	return
}

func (c CandleHandler) Patch(w http.ResponseWriter, r *http.Request) {
	util.WriteError(w, http.StatusNotImplemented)
	return
}

func (c CandleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	util.WriteError(w, http.StatusNotImplemented)
	return
}

func (c CandleHandler) Options(w http.ResponseWriter, r *http.Request) {
	return
}

func (c CandleHandler) Handle(r *http.ServeMux) {
	r.HandleFunc("GET /candle/{ticker}", c.Get)
	r.HandleFunc("POST /candle/{ticker}", c.Post)
	r.HandleFunc("PUT /candle/{ticker}", c.Put)
	r.HandleFunc("PATCH /candle/{ticker}", c.Patch)
	r.HandleFunc("DELETE /candle/{ticker}", c.Delete)
	r.HandleFunc("OPTIONS /candle/{ticker}", c.Options)
}
