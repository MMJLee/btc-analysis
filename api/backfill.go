package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/database"
)

type BackfillHandler struct {
	pool      database.DBPool
	tickerMap map[string]chan bool
	mut       *sync.Mutex
}

func NewBackfillHandler(pool database.DBPool, tickerMap map[string]chan bool, mut *sync.Mutex) *BackfillHandler {
	return &BackfillHandler{pool, tickerMap, mut}
}

func (b *BackfillHandler) requireAuth() bool {
	return true
}

func (b *BackfillHandler) get(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")

	b.mut.Lock()
	_, exists := b.tickerMap[ticker]
	b.mut.Unlock()

	message := fmt.Sprintf("Not backfilling %s", ticker)
	if exists {
		message = fmt.Sprintf("Currently backfilling %s", ticker)
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (b *BackfillHandler) post(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	queryParams := r.URL.Query()
	start := queryParams.Get("start")
	end := queryParams.Get("end")
	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}
	b.mut.Lock()
	if len(b.tickerMap) > 1 {
		b.mut.Unlock()
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	if _, exists := b.tickerMap[ticker]; exists {
		b.mut.Unlock()
		writeError(w, http.StatusConflict)
		return
	}
	stopChan := make(chan bool)
	b.tickerMap[ticker] = stopChan
	b.mut.Unlock()

	go func() {
		client.BackfillTicker(ticker, startInt, endInt, stopChan)
		b.mut.Lock()
		delete(b.tickerMap, ticker)
		b.mut.Unlock()
	}()
	jsonData, err := json.Marshal(fmt.Sprintf("Now backfilling %s", ticker))
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (b *BackfillHandler) put(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (b *BackfillHandler) patch(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (b *BackfillHandler) delete(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	message := fmt.Sprintf("Not backfilling %s", ticker)

	b.mut.Lock()
	stopChan, exists := b.tickerMap[ticker]
	if exists {
		stopChan <- true
		delete(b.tickerMap, ticker)
		message = fmt.Sprintf("Stopped backfilling %s", ticker)
	}
	b.mut.Unlock()

	jsonData, err := json.Marshal(message)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (b *BackfillHandler) options(w http.ResponseWriter, r *http.Request) {
}

func (b *BackfillHandler) handle(r *http.ServeMux) {
	r.HandleFunc("GET /backfill/{ticker}", b.get)
	r.HandleFunc("POST /backfill/{ticker}", b.post)
	r.HandleFunc("PUT /backfill/{ticker}", b.put)
	r.HandleFunc("PATCH /backfill/{ticker}", b.patch)
	r.HandleFunc("DELETE /backfill/{ticker}", b.delete)
	r.HandleFunc("OPTIONS /backfill/{ticker}", b.options)
}
