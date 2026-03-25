// Command migrate is a standalone CLI for running database migrations.
//
// Usage:
//
//	migrate -cmd up                  # apply all pending migrations
//	migrate -cmd down                # roll back all applied migrations
//	migrate -cmd steps -steps 2      # apply 2 migrations forward
//	migrate -cmd steps -steps -1     # roll back 1 migration
//	migrate -cmd version             # show current version
//
// Connection parameters are read from the environment (same .env as the main
// server), so the command can be run locally or in a CI/CD pipeline:
//
//	GONE_DB_HOST=localhost GONE_DB_PORT=5432 ... migrate -cmd up
//
// or simply:
//
//	make migrate-up   (see Makefile)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	dbmigrate "starter/internal/db/migrate"
)

func main() {
	cmd := flag.String("cmd", "up", "Migration command: up | down | steps | version")
	steps := flag.Int("steps", 0, "Steps to migrate (+N forward, -N backward) — only for -cmd steps")
	flag.Parse()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("GONE_DB_HOST"),
		os.Getenv("GONE_DB_USERNAME"),
		os.Getenv("GONE_DB_PASSWORD"),
		os.Getenv("GONE_DB_DATABASE"),
		os.Getenv("GONE_DB_PORT"),
	)

	if err := dbmigrate.Run(dsn, *cmd, *steps); err != nil {
		log.Fatalf("migration error: %v", err)
	}
}
