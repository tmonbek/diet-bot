package postgres

import (
	"context"
	"log"
	"log/slog"
	"os"

	"diet-bot/internal/domain/repository"
	"diet-bot/internal/infrastructure/persistence"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var UserRepo repository.UserRepository

func InitDB() {
	connStr := os.Getenv("CONNECTION_STRING")
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	UserRepo = persistence.NewPostgresUserRepository(Pool)
	slog.Info("Connected to database")
}

func CloseDB() {
	Pool.Close()
}
