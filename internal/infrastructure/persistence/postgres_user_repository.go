package persistence

import (
	"context"
	"log"

	"diet-bot/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, userID int64) (int, error) {
	var newID int
	err := r.pool.QueryRow(ctx,
		"INSERT INTO users (user_id) VALUES ($1) RETURNING id",
		userID,
	).Scan(&newID)

	if err != nil {
		log.Printf("Insert error: %v", err)
		return 0, err
	}

	return newID, nil
}

func (r *PostgresUserRepository) Exists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)",
		userID,
	).Scan(&exists)

	if err != nil {
		log.Printf("Query error: %v", err)
		return false, err
	}

	return exists, nil
}

func (r *PostgresUserRepository) GetAllChats(ctx context.Context) ([]int64, error) {
	var chatIDs []int64

	rows, err := r.pool.Query(ctx, "SELECT user_id FROM users")
	if err != nil {
		log.Printf("Failed to query users: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		return nil, err
	}

	return chatIDs, nil
}
