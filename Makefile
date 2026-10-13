BACKEND := $(CURDIR)/backend
# sqlc runs in a container: building it locally needs cgo + gcc.
SQLC    := docker run --rm -u $(shell id -u):$(shell id -g) -v $(BACKEND):/src -w /src sqlc/sqlc:1.29.0

.PHONY: db run up down logs sqlc tidy test frontend

db:          ## Start only Postgres (localhost:5432)
	docker compose up -d db

run: db      ## Run the backend locally with Go (http://localhost:8080)
	cd backend && go run ./cmd/server

up:          ## Start Postgres + backend, both in Docker
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f backend

sqlc:        ## Regenerate backend/internal/store from db/queries
	$(SQLC) generate

tidy:
	cd backend && go mod tidy

test:
	cd backend && go vet ./... && go test ./...

frontend:    ## Start the Vite dev server (http://localhost:5173)
	cd frontend && npm run dev
