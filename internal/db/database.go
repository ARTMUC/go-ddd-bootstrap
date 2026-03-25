package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *gorm.DB
}

type service struct {
	db *gorm.DB
}

var (
	dbName   = os.Getenv("GONE_DB_DATABASE")
	password = os.Getenv("GONE_DB_PASSWORD")
	username = os.Getenv("GONE_DB_USERNAME")
	port     = os.Getenv("GONE_DB_PORT")
	host     = os.Getenv("GONE_DB_HOST")
)

// New connects to the database and returns a Service.
// Migrations are NOT run automatically – use `cmd/migrate` to apply them.
func New() Service {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, dbName, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &service{db: db}
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, _ := s.db.DB()
	if err := sqlDB.PingContext(ctx); err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}
	return map[string]string{"status": "up"}
}

func (s *service) Close() error {
	sqlDB, _ := s.db.DB()
	return sqlDB.Close()
}
