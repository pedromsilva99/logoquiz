package api

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"

	"github.com/pedro/logoquiz/backend/internal/store"
)

var imageExts = []string{".svg", ".png", ".webp", ".jpg", ".jpeg"}

// logoImage serves the hidden (incomplete) version of a logo until the player
// solves it, and the full version afterwards. The real filename never reaches
// the client, so it can't give the answer away.
func (s *Server) logoImage(w http.ResponseWriter, r *http.Request) {
	logoID, ok := parseLogoID(w, r)
	if !ok {
		return
	}

	ctx := r.Context()
	logo, err := s.q.GetLogo(ctx, logoID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "logo not found")
		return
	}
	if err != nil {
		internalError(w, r, err)
		return
	}

	solved, err := s.q.IsLogoSolved(ctx, store.IsLogoSolvedParams{PlayerID: playerID(ctx), LogoID: logoID})
	if err != nil {
		internalError(w, r, err)
		return
	}

	variant := "hidden"
	if solved {
		variant = "full"
	}

	// The same URL returns different images depending on the player's progress.
	w.Header().Set("Cache-Control", "private, no-cache")

	for _, ext := range imageExts {
		path := filepath.Join(s.assetsDir, variant, logo.ImageKey+ext)
		if _, err := os.Stat(path); err == nil {
			http.ServeFile(w, r, path)
			return
		}
	}

	// No image on disk yet: serve a placeholder so the game is playable without assets.
	label := "?"
	if solved {
		label = logo.Name
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprintf(w, placeholderSVG, html.EscapeString(label))
}

const placeholderSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200">
<rect width="200" height="200" rx="24" fill="#e8e4dc"/>
<text x="100" y="100" font-family="system-ui, sans-serif" font-size="22" font-weight="700"
      fill="#5b544a" text-anchor="middle" dominant-baseline="middle">%s</text>
</svg>`
