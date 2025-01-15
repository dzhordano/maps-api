package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		fmt.Println("connString", connString)

		return nil, err
	}

	timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(timeout); err != nil {
		return nil, err
	}

	return pool, nil
}

func Close(pool *pgxpool.Pool) {
	pool.Close()
}
