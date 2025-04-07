package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

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
}

func CloseDB() {
	Pool.Close()
}

func GetAllChats() []int64 {
	var chatIDs []int64

	rows, err := Pool.Query(context.Background(), "SELECT user_id FROM users")
	if err != nil {
		log.Printf("Failed to query users: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue // Пропускаем ошибочную строку, не прерываем всю функцию
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
	}

	return chatIDs
}
