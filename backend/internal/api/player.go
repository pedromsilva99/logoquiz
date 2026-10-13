package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Players are anonymous for now: identified only by a UUID kept in a cookie.
const playerCookie = "player_id"

type playerKey struct{}

func (s *Server) withPlayer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := s.resolvePlayer(w, r)
		if err != nil {
			internalError(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), playerKey{}, id)))
	})
}

// resolvePlayer returns the player from the cookie, creating a new one when the
// cookie is missing, malformed or points to a player that no longer exists.
func (s *Server) resolvePlayer(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	if c, err := r.Cookie(playerCookie); err == nil {
		if id, err := uuid.Parse(c.Value); err == nil {
			exists, err := s.q.PlayerExists(r.Context(), id)
			if err != nil {
				return uuid.Nil, err
			}
			if exists {
				return id, nil
			}
		}
	}

	id, err := s.q.CreatePlayer(r.Context())
	if err != nil {
		return uuid.Nil, err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     playerCookie,
		Value:    id.String(),
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return id, nil
}

func playerID(ctx context.Context) uuid.UUID {
	return ctx.Value(playerKey{}).(uuid.UUID)
}
