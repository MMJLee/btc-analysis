package api

import (
	"net/http"

	"github.com/mmjlee/btc-analysis/internal/util"
)

func GetServer(c CandleHandler) http.Server {
	router := http.NewServeMux()
	router.HandleFunc("GET /candle/{ticker}", c.GetCandles)
	router.HandleFunc("POST /candle/{ticker}", c.GetCandles)
	router.HandleFunc("OPTIONS /candle/{ticker}", c.Options)
	admin_router := http.NewServeMux()
	admin_router.HandleFunc("PUT /candle/{ticker}", c.GetCandles)
	admin_router.HandleFunc("PATCH /candle/{ticker}", c.GetCandles)
	admin_router.HandleFunc("DELETE /candle/{ticker}", c.GetCandles)
	router.Handle("/", util.AuthMiddleware(admin_router))
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))
	middlewares := util.CreateStack(
		util.GzipMiddleware,
		util.CORSMiddleware,
		util.ErrorMiddleware,
		util.LoggingMiddleware,
	)

	return http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: middlewares(v1),
	}
}
