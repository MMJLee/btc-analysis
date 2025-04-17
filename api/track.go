package api

import (
	"encoding/json"
	"fmt"
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

func (t *TrackHandler) requireAuth() bool {
	return true
}

func (t *TrackHandler) get(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *TrackHandler) post(w http.ResponseWriter, r *http.Request) {
	ticker := r.PathValue("ticker")

	t.mut.Lock()
	if len(t.trackMap) > 1 {
		t.mut.Unlock()
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	if _, exists := t.trackMap[ticker]; exists {
		t.mut.Unlock()
		writeError(w, http.StatusConflict)
		return
	}
	stopChan := make(chan bool)
	t.trackMap[ticker] = stopChan
	t.mut.Unlock()

	go client.TrackTicker(ticker, stopChan)
	jsonData, err := json.Marshal(fmt.Sprintf("Now tracking %s", ticker))
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *TrackHandler) put(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *TrackHandler) patch(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *TrackHandler) delete(w http.ResponseWriter, r *http.Request) {
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
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *TrackHandler) options(w http.ResponseWriter, r *http.Request) {
}

func (t *TrackHandler) handle(r *http.ServeMux) {
	r.HandleFunc("GET /track/{ticker}", t.get)
	r.HandleFunc("POST /track/{ticker}", t.post)
	r.HandleFunc("PUT /track/{ticker}", t.put)
	r.HandleFunc("PATCH /track/{ticker}", t.patch)
	r.HandleFunc("DELETE /track/{ticker}", t.delete)
	r.HandleFunc("OPTIONS /track/{ticker}", t.options)
}
