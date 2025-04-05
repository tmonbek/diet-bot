package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	connStr := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	var err error
	Pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to create pool:", err)
	}
}

func CloseDB() {
	Pool.Close()
}
