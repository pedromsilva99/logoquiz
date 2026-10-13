package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/pedro/logoquiz/backend/internal/game"
	"github.com/pedro/logoquiz/backend/internal/store"
)

type logoDTO struct {
	ID       int32  `json:"id"`
	Position int32  `json:"position"`
	Solved   bool   `json:"solved"`
	Name     string `json:"name,omitempty"` // only revealed once solved
}

type levelDTO struct {
	Number int32     `json:"number"`
	Title  string    `json:"title"`
	Logos  []logoDTO `json:"logos"`
}

func (s *Server) getLevel(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.Atoi(r.PathValue("number"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid level number")
		return
	}

	ctx := r.Context()
	level, err := s.q.GetLevelByNumber(ctx, int32(number))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "level not found")
		return
	}
	if err != nil {
		internalError(w, r, err)
		return
	}

	rows, err := s.q.ListLevelLogosForPlayer(ctx, store.ListLevelLogosForPlayerParams{
		LevelID:  level.ID,
		PlayerID: playerID(ctx),
	})
	if err != nil {
		internalError(w, r, err)
		return
	}

	logos := make([]logoDTO, 0, len(rows))
	for _, row := range rows {
		l := logoDTO{ID: row.ID, Position: row.Position, Solved: row.Solved}
		if row.Solved {
			l.Name = row.Name
		}
		logos = append(logos, l)
	}
	writeJSON(w, http.StatusOK, levelDTO{Number: level.Number, Title: level.Title, Logos: logos})
}

type guessRequest struct {
	Answer string `json:"answer"`
}

type guessResponse struct {
	Correct bool   `json:"correct"`
	Name    string `json:"name,omitempty"`
}

func (s *Server) guessLogo(w http.ResponseWriter, r *http.Request) {
	logoID, ok := parseLogoID(w, r)
	if !ok {
		return
	}

	var req guessRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	ctx := r.Context()
	player := playerID(ctx)
	logo, err := s.q.GetLogo(ctx, logoID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "logo not found")
		return
	}
	if err != nil {
		internalError(w, r, err)
		return
	}

	solved, err := s.q.IsLogoSolved(ctx, store.IsLogoSolvedParams{PlayerID: player, LogoID: logoID})
	if err != nil {
		internalError(w, r, err)
		return
	}
	if solved {
		writeJSON(w, http.StatusOK, guessResponse{Correct: true, Name: logo.Name})
		return
	}

	aliases, err := s.q.ListLogoAliases(ctx, logoID)
	if err != nil {
		internalError(w, r, err)
		return
	}
	correct := game.Matches(req.Answer, append([]string{logo.Name}, aliases...))

	if err := s.q.RecordAttempt(ctx, store.RecordAttemptParams{
		PlayerID: player,
		LogoID:   logoID,
		Correct:  correct,
	}); err != nil {
		internalError(w, r, err)
		return
	}

	res := guessResponse{Correct: correct}
	if correct {
		res.Name = logo.Name
	}
	writeJSON(w, http.StatusOK, res)
}

func parseLogoID(w http.ResponseWriter, r *http.Request) (int32, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid logo id")
		return 0, false
	}
	return int32(id), true
}
