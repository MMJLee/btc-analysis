package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/mmjlee/btc-analysis/internal/client"
	"github.com/mmjlee/btc-analysis/internal/database"
)

type TrackHandler struct {
	pool     database.DBPool
	trackMap map[string]chan bool
	mut      *sync.Mutex
}

func NewTrackHandler(pool database.DBPool, trackMap map[string]chan bool, mut *sync.Mutex) *TrackHandler {
	return &TrackHandler{pool: pool, trackMap: trackMap, mut: mut}
}

func (t *TrackHandler) Get(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")

	t.mut.Lock()
	_, exists := t.trackMap[ticker]
	t.mut.Unlock()

	message := fmt.Sprintf("Not tracking %s", ticker)
	if exists {
		message = fmt.Sprintf("Currently tracking %s", ticker)
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Panicf("Track-Get-%w", err)
	}
	w.Write(jsonData)
}

func (t *TrackHandler) Post(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")

	t.mut.Lock()
	if len(t.trackMap) > 1 {
		t.mut.Unlock()
		WriteError(w, http.StatusServiceUnavailable)
		return
	}
	if _, exists := t.trackMap[ticker]; exists {
		t.mut.Unlock()
		WriteError(w, http.StatusConflict)
		return
	}
	stopChan := make(chan bool)
	t.trackMap[ticker] = stopChan
	t.mut.Unlock()

	go client.TrackTicker(ticker, stopChan)
	jsonData, err := json.Marshal(fmt.Sprintf("Now tracking %s", ticker))
	if err != nil {
		log.Panicf("Track-Post-%w", err)
	}
	w.Write(jsonData)
}

func (t *TrackHandler) Put(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotImplemented)
	return
}

func (t *TrackHandler) Patch(w http.ResponseWriter, r *http.Request) {
	WriteError(w, http.StatusNotImplemented)
	return
}

func (t *TrackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")
	message := fmt.Sprintf("Not tracking %s", ticker)

	t.mut.Lock()
	stopChan, exists := t.trackMap[ticker]
	if exists {
		stopChan <- true
		delete(t.trackMap, ticker)
		message = fmt.Sprintf("Stopped tracking %s", ticker)
	}
	t.mut.Unlock()

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Panicf("Track-Delete-%w", err)
	}
	w.Write(jsonData)
}

func (t *TrackHandler) Options(w http.ResponseWriter, r *http.Request) {
	return
}

func (t *TrackHandler) Handle(r *http.ServeMux) {
	r.HandleFunc("GET /track/{ticker}", t.Get)
	r.HandleFunc("POST /track/{ticker}", t.Post)
	r.HandleFunc("PUT /track/{ticker}", t.Put)
	r.HandleFunc("PATCH /track/{ticker}", t.Patch)
	r.HandleFunc("DELETE /track/{ticker}", t.Delete)
	r.HandleFunc("OPTIONS /track/{ticker}", t.Options)
}

func (t *TrackHandler) RequireAuth() bool {
	return true
}
