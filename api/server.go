package api

import (
	"net/http"
)

type Handler interface {
	requireAuth() bool
	get(w http.ResponseWriter, r *http.Request)
	post(w http.ResponseWriter, r *http.Request)
	put(w http.ResponseWriter, r *http.Request)
	patch(w http.ResponseWriter, r *http.Request)
	delete(w http.ResponseWriter, r *http.Request)
	options(w http.ResponseWriter, r *http.Request)
	handle(r *http.ServeMux)
}

func GetServer(handlers ...Handler) http.Server {
	baseMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	for _, h := range handlers {
		if h.requireAuth() {
			h.handle(adminMux)
		} else {
			h.handle(baseMux)
		}
	}

	baseMux.Handle("/", authMiddleware(adminMux))
	middledMux := applyMiddlewares(baseMux)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", middledMux))
	return http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: v1,
	}
}
