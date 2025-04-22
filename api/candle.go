package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mmjlee/btc-analysis/internal/database"
)

type CandleHandler struct {
	pool database.DBPool
}

func NewCandleHandler(pool database.DBPool) *CandleHandler {
	return &CandleHandler{pool}
}

func (c *CandleHandler) requireAuth() bool {
	return false
}

func (c *CandleHandler) get(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	queryParams := r.URL.Query()
	start := queryParams.Get("start")
	end := queryParams.Get("end")
	limit := queryParams.Get("limit")
	offset := queryParams.Get("offset")
	if ticker == "" || start == "" || end == "" || limit == "" || offset == "" {
		writeError(w, http.StatusBadRequest)
		return
	}
	missing, _ := strconv.ParseBool(queryParams.Get("missing"))
	if missing { //make sure start epoch is truncated to the minute
		startInt, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest)
			return
		}
		start = strconv.FormatInt(time.Unix(startInt, 0).Truncate(time.Minute).Unix(), 10)
	}
	candles, err := c.pool.GetCandles(r.Context(), ticker, start, end, limit, offset, missing)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(candles)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (c *CandleHandler) post(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (c *CandleHandler) put(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (c *CandleHandler) patch(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (c *CandleHandler) delete(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (c *CandleHandler) options(w http.ResponseWriter, r *http.Request) {
}

func (c *CandleHandler) handle(r *http.ServeMux) {
	r.HandleFunc("GET /candle/{ticker}", c.get)
	r.HandleFunc("POST /candle/{ticker}", c.post)
	r.HandleFunc("PUT /candle/{ticker}", c.put)
	r.HandleFunc("PATCH /candle/{ticker}", c.patch)
	r.HandleFunc("DELETE /candle/{ticker}", c.delete)
	r.HandleFunc("OPTIONS /candle/{ticker}", c.options)
}
