package api

import (
	"encoding/json"
	"fmt"
	"log"
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
	return &BackfillHandler{pool: pool, tickerMap: tickerMap, mut: mut}
}

func (b *BackfillHandler) Get(w http.ResponseWriter, r *http.Request) {
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
		log.Panicf("Backfill-Get-%v", err)
	}
	w.Write(jsonData)
}

func (b *BackfillHandler) Post(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	queryParams := r.URL.Query()
	start := queryParams.Get("start")
	end := queryParams.Get("end")
	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest)
		return
	}
	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest)
		return
	}
	b.mut.Lock()
	if len(b.tickerMap) > 1 {
		b.mut.Unlock()
		WriteError(w, http.StatusServiceUnavailable)
		return
	}
	if _, exists := b.tickerMap[ticker]; exists {
		b.mut.Unlock()
		WriteError(w, http.StatusConflict)
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
		log.Panicf("Backfill-Post-%w", err)
	}
	w.Write(jsonData)
}

func (b *BackfillHandler) Put(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotImplemented)
	return
}

func (b *BackfillHandler) Patch(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotImplemented)
	return
}

func (b *BackfillHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		log.Panicf("Backfill-Delete-%w", err)

	}
	w.Write(jsonData)
}

func (b *BackfillHandler) Options(w http.ResponseWriter, r *http.Request) {
	return
}

func (b *BackfillHandler) Handle(r *http.ServeMux) {
	r.HandleFunc("GET /backfill/{ticker}", b.Get)
	r.HandleFunc("POST /backfill/{ticker}", b.Post)
	r.HandleFunc("PUT /backfill/{ticker}", b.Put)
	r.HandleFunc("PATCH /backfill/{ticker}", b.Patch)
	r.HandleFunc("DELETE /backfill/{ticker}", b.Delete)
	r.HandleFunc("OPTIONS /backfill/{ticker}", b.Options)
}

func (b *BackfillHandler) RequireAuth() bool {
	return true
}
