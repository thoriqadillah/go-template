package queue

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

func Start(ctx context.Context, pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	dbpool := riverpgxv5.New(pool)
	client, err := river.NewClient(dbpool, &river.Config{
		Workers: workers,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: 100,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return client, err
}
