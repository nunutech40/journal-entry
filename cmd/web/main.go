package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"journal-entry/internal/account"
	"journal-entry/internal/journal"
	"journal-entry/internal/report"
	"journal-entry/internal/server"
)

func main() {
	// Load .env file (ignore error if not found — production uses real env vars)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Init database connection pool
	pool, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("✅ Database connected")

	// Parse all HTML templates at startup (fail fast)
	templates := server.ParseTemplates()

	// Wire dependencies: Repository → Service → Handler
	accountRepo := account.NewRepository(pool)
	accountSvc := account.NewService(accountRepo)
	accountHandler := account.NewHandler(accountSvc, templates)

	journalRepo := journal.NewRepository(pool)
	journalSvc := journal.NewService(journalRepo, accountRepo)
	journalHandler := journal.NewHandler(journalSvc, accountSvc, templates)

	reportRepo := report.NewRepository(pool)
	reportSvc := report.NewService(reportRepo, accountRepo)
	reportHandler := report.NewHandler(reportSvc, accountSvc, templates)

	// Create router with all routes
	router := server.NewRouter(templates, accountHandler, journalHandler, reportHandler)

	// Get port from env
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server with graceful shutdown
	srv := server.New(router, port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// initDB creates a PostgreSQL connection pool using environment variables.
func initDB() (*pgxpool.Pool, error) {
	host := envOrDefault("DB_HOST", "localhost")
	port := envOrDefault("DB_PORT", "5432")
	user := envOrDefault("DB_USER", "postgres")
	pass := envOrDefault("DB_PASS", "postgres")
	name := envOrDefault("DB_NAME", "journal_entry")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}

func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
