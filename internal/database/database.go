package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect 创建 PostgreSQL 连接池；ctx 为启动上下文，databaseURL 为 PostgreSQL 连接串。
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
