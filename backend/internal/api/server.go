// Package api exposes the game over HTTP.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pedro/logoquiz/backend/internal/store"
)

type Server struct {
	q         *store.Queries
	assetsDir string
}

func NewServer(q *store.Queries, assetsDir string) *Server {
	return &Server{q: q, assetsDir: assetsDir}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/levels/{number}", s.getLevel)
	mux.HandleFunc("POST /api/logos/{id}/guess", s.guessLogo)
	mux.HandleFunc("GET /api/logos/{id}/image", s.logoImage)
	return s.withPlayer(mux)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
	writeError(w, http.StatusInternalServerError, "internal error")
}
