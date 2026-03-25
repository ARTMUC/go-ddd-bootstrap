// Package migrate provides a reusable migration runner backed by golang-migrate.
// SQL files are embedded directly into the binary so no external migration files
// need to be present at runtime.
package migrate

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed *.sql
var migrationsFS embed.FS

// Run opens a postgres connection with the given DSN and runs the requested
// migration command.
//
// Supported commands:
//
//	up      – apply all pending migrations
//	down    – roll back all applied migrations
//	steps N – apply (+N) or roll back (-N) exactly N migrations
//	version – print the current migration version and dirty flag, then exit
func Run(dsn, command string, steps int) error {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer sqlDB.Close()

	src, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
		log.Println("migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate down: %w", err)
		}
		log.Println("migrations rolled back successfully")

	case "steps":
		if steps == 0 {
			return fmt.Errorf("steps command requires a non-zero -steps value")
		}
		if err := m.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate steps %d: %w", steps, err)
		}
		log.Printf("migrated %d step(s) successfully", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("get version: %w", err)
		}
		log.Printf("migration version: %d, dirty: %v", version, dirty)

	default:
		return fmt.Errorf("unknown command %q – supported: up, down, steps, version", command)
	}

	return nil
}
