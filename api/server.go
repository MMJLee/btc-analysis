package api

import (
	"net/http"

	"github.com/mmjlee/btc-analysis/internal/util"
)

type Handler interface {
	Handle(r *http.ServeMux)
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Options(w http.ResponseWriter, r *http.Request)
}

func GetServer(handlers ...Handler) http.Server {
	baseHandler := http.NewServeMux()
	for _, h := range handlers {
		h.Handle(baseHandler)
	}
	handler := util.ApplyMiddlewares(baseHandler)
	// admin_router := http.NewServeMux()
	// admin_router.HandleFunc("PUT /candle/{ticker}", c.GetCandles)
	// router.Handle("/", util.AuthMiddleware(admin_router))
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", handler))
	return http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: v1,
	}
}
