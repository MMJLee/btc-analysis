package api

import (
	"fmt"
	"net/http"
)

func HandleCandle(w http.ResponseWriter, r *http.Request) {
	product_id := r.PathValue("product_id")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	w.Write([]byte(fmt.Sprintf("%v, %v, %v", product_id, start, end)))
}
