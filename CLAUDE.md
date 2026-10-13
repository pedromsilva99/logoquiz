# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A "Logo Quiz" web game focused on Portuguese brands. Go backend + PostgreSQL, React (Vite + TypeScript) frontend. User-facing text is in European Portuguese (pt-PT).

## Commands

Run from the repo root (inside WSL):

```sh
make run        # Postgres in Docker + backend with local Go on :8080 (applies migrations on startup)
make frontend   # Vite dev server on :5173, proxies /api -> :8080
make test       # go vet + go test for the backend
make sqlc       # regenerate backend/internal/store after editing db/queries or db/migrations
make tidy       # go mod tidy
make up / down  # full stack (db + backend) in Docker
```

Single Go test: `cd backend && go test ./internal/game -run TestMatches`
Frontend type-check + build: `cd frontend && npm run build`

Environment notes:
- Go lives in `/usr/local/go/bin`, which is only on PATH in interactive shells. When invoking via `wsl ... bash`, use `bash -lic`.
- The Docker `buildx` plugin on this machine is a broken symlink to Docker Desktop, so `make up` (which builds the backend image) fails. Use `make run`, or `DOCKER_BUILDKIT=0 docker build -t logoquiz-backend ./backend && docker compose up -d --no-build`.
- `sqlc` runs via the `sqlc/sqlc` Docker image, not a local binary.

## Architecture

**Request flow.** `cmd/server/main.go` opens a pgx pool, runs goose migrations embedded from `backend/db` (`db.Migrations`), then serves `api.Server.Routes()`. Every route is wrapped by `withPlayer` (`internal/api/player.go`). It resolves an anonymous player from the `player_id` cookie, creating a row in `players` and setting the cookie if needed, and puts the UUID in the request context (`playerID(ctx)`). There are no accounts yet.

**Data access is sqlc-generated.** Write SQL in `backend/db/queries/*.sql` and the schema in `backend/db/migrations/*.sql` (goose format, `-- +goose Up/Down`), then run `make sqlc`. Never hand-edit `backend/internal/store/`. The `uuid` columns map to `github.com/google/uuid.UUID` (see `sqlc.yaml`). Seed data (levels, logos, aliases) is also added as migrations, e.g. `00002_seed_level1.sql`.

**Never leak the answer to the client.** This is the core invariant:
- `GET /api/levels/{n}` includes a logo's `name` only once that player has solved it.
- Images are served only via `GET /api/logos/{id}/image`. The server looks up `image_key` and the player's progress, then serves `assets/logos/hidden/<key>.<ext>` before the logo is solved and `assets/logos/full/<key>.<ext>` after. It falls back to a placeholder SVG when the file is missing. `image_key` and filenames must never reach the frontend.
- Answers are checked server-side in `POST /api/logos/{id}/guess` with `game.Matches` (`internal/game/answer.go`). That function normalizes input (lowercase, strips accents and non-alphanumerics) and allows Levenshtein typos scaled by length, so names of 4 characters or fewer must be exact. Accepted answers are the logo name plus the rows in `logo_aliases`.
- Progress lives in `progress (player_id, logo_id, attempts, solved_at)`. `RecordAttempt` is an upsert that never clears `solved_at`.

**Assets.** Images are dropped manually into `backend/assets/logos/{hidden,full}/`, named by `image_key` (list in `backend/assets/logos/README.md`). They are read from disk per request, and in Docker the folder is mounted read-only, so no restart is needed.

**Frontend.** `src/api.ts` holds the API types and fetchers. Server state goes through TanStack Query, with the level under query key `['level', n]`. On a correct guess, `App.tsx` updates that cache with `setQueryData` instead of refetching. The image URL carries `?v=hidden|full` only to bust the browser cache when a logo flips, because the server decides which variant to send. In dev, Vite proxies `/api`, so the player cookie works same-origin.
