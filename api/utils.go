package api

import (
	"compress/gzip"
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const AUTH_USER = "middleware.auth.user"

type Middleware func(http.Handler) http.Handler

func createStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range middlewares {
			next = middleware(next)
		}
		return next
	}
}

func applyMiddlewares(handler *http.ServeMux) http.Handler {
	middlewares := createStack(
		gzipMiddleware,
		corsMiddleware,
		errorMiddleware,
		loggingMiddleware,
	)
	return middlewares(handler)
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				writeError(w, http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		log.Println(wrapped.statusCode, r.Method, r.URL, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://mjlee.dev")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Authorization, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (grw *gzipResponseWriter) Write(p []byte) (n int, err error) {
	return grw.Writer.Write(p)
}
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gz := gzip.NewWriter(w)
			defer gz.Close()
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")
			w = &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
			}
		}
		next.ServeHTTP(w, r)
	})
}

// WIP
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		token, err := base64.StdEncoding.DecodeString(authorization)
		if err != nil || authorization == "" {
			writeError(w, http.StatusUnauthorized)
			return
		}
		userID := string(token)
		ctx := context.WithValue(r.Context(), AUTH_USER, userID)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func writeError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

func hashPassword(password string) (string, error) {
	// salting done by bcrypt
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
