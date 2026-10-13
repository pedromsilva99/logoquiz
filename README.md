# logoquiz

Jogo estilo *Logo Quiz* com marcas portuguesas.

- **backend/** — Go (`net/http`, pgx, sqlc, goose), PostgreSQL
- **frontend/** — React + Vite + TypeScript, TanStack Query

## Correr

Requisitos: Go 1.26+, Node e Docker (para o Postgres e o sqlc).

```sh
make run         # Postgres em Docker + backend com Go em http://localhost:8080 (migrações correm no arranque)
make frontend    # (cd frontend && npm install na primeira vez) → http://localhost:5173
```

Alternativa: `make up` corre também o backend em Docker.

## Imagens

Ver [backend/assets/logos/README.md](backend/assets/logos/README.md). As imagens são lidas do disco
a cada pedido, por isso basta copiá-las para lá — não é preciso reiniciar nada.

## Desenvolvimento

```sh
make sqlc        # depois de alterar backend/db/queries/*.sql ou as migrações
make tidy        # go mod tidy
make test        # go vet + go test
```

Novas migrações vão em `backend/db/migrations/` (formato goose, `000NN_nome.sql`).
