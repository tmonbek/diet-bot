package repository

import "context"

type UserRepository interface {
	Create(ctx context.Context, userID int64) (int, error)
	Exists(ctx context.Context, userID int64) (bool, error)
	GetAllChats(ctx context.Context) ([]int64, error)
}
