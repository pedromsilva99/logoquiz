package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/pedro/logoquiz/backend/db"
	"github.com/pedro/logoquiz/backend/internal/api"
	"github.com/pedro/logoquiz/backend/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	dsn := env("DATABASE_URL", "postgres://logoquiz:logoquiz@localhost:5432/logoquiz?sslmode=disable")
	addr := env("ADDR", ":8080")
	assetsDir := env("ASSETS_DIR", "assets/logos")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := migrate(pool); err != nil {
		return err
	}

	srv := api.NewServer(store.New(pool), assetsDir)
	slog.Info("listening", "addr", addr)
	return http.ListenAndServe(addr, srv.Routes())
}

func migrate(pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	goose.SetBaseFS(db.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, "migrations")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
