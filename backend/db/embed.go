// Package db holds the SQL migrations and queries. Migrations are embedded so the
// server binary can apply them on startup.
package db

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
