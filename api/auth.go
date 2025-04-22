package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mmjlee/btc-analysis/internal/database"
)

type AuthHandler struct {
	pool database.DBPool
	rdb  database.RedisClient
}

func NewAuthHandler(pool database.DBPool, rdb database.RedisClient) *AuthHandler {
	return &AuthHandler{pool, rdb}
}

func (t *AuthHandler) get(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *AuthHandler) post(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *AuthHandler) put(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *AuthHandler) patch(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *AuthHandler) delete(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented)
}

func (t *AuthHandler) options(w http.ResponseWriter, r *http.Request) {
}

func (t *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		writeError(w, http.StatusNotAcceptable)
		return
	}
	if _, err := t.pool.GetUser(r.Context(), username); err == nil {
		writeError(w, http.StatusNotAcceptable)
		return
	}
	hashedPassword, err := hashPassword(password)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	err = t.pool.CreateUser(r.Context(), username, hashedPassword)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Registered %s", username)
	jsonData, err := json.Marshal(message)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	timeout := time.Duration(24) * time.Hour
	user, err := t.pool.GetUser(r.Context(), username)
	if err != nil || !checkHash(user.Password, password) {
		writeError(w, http.StatusUnauthorized)
		return
	}
	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	userSession, err := json.Marshal(SessionData{username, csrfToken, sessionToken})
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	err = t.rdb.Set(r.Context(), username, userSession, timeout).Err()
	if err != nil {
		writeError(w, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_TOKEN,
		Value:    sessionToken,
		Expires:  time.Now().Add(timeout),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CSRF_TOKEN,
		Value:    csrfToken,
		Expires:  time.Now().Add(timeout),
		HttpOnly: false,
	})

	message := "Logged in"
	jsonData, err := json.Marshal(message)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_TOKEN,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     CSRF_TOKEN,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	t.rdb.Del(r.Context(), username)

	message := "Logged out"
	jsonData, err := json.Marshal(message)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (t *AuthHandler) handle(r *http.ServeMux) {
	r.HandleFunc("GET /auth", t.get)
	r.HandleFunc("POST /auth", t.post)
	r.HandleFunc("PUT /auth", t.put)
	r.HandleFunc("PATCH /auth", t.patch)
	r.HandleFunc("DELETE /auth", t.delete)
	r.HandleFunc("OPTIONS /auth", t.options)
	r.HandleFunc("POST /auth/register", t.register)
	r.HandleFunc("POST /auth/login", t.login)
	r.HandleFunc("GET /auth/logout", t.logout)
}

func (t *AuthHandler) requireAuth() bool {
	return false
}
