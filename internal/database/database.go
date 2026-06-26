package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Use the exact connection string format from Supabase
	connStr := fmt.Sprintf("postgresql://postgres.hwzvpczxhfddgwcxqxxp:%s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres",
		password)

	log.Printf("Connecting to database using Supabase connection string\n")

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Increase timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Connection string:", connStr)
	log.Println("Successfully connected to Supabase database")

	// Configure connection pool with more conservative settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	dbInstance = &service{db: db}
	return dbInstance
}

// Health returns the health status of the database connection
func (s *service) Health() map[string]string {
	health := make(map[string]string)

	// Increase health check timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		health["status"] = "unhealthy"
		health["message"] = err.Error()
	} else {
		health["status"] = "healthy"
	}

	return health
}

// Close closes the database connection
func (s *service) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
