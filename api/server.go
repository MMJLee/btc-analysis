package api

import (
	"net/http"
)

type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Options(w http.ResponseWriter, r *http.Request)
	Handle(r *http.ServeMux)
	RequireAuth() bool
}

func GetServer(handlers ...Handler) http.Server {
	baseMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	for _, h := range handlers {
		if h.RequireAuth() {
			h.Handle(adminMux)
		} else {
			h.Handle(baseMux)
		}
	}

	baseMux.Handle("/", AuthMiddleware(adminMux))
	middledMux := ApplyMiddlewares(baseMux)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", middledMux))
	return http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: v1,
	}
}
