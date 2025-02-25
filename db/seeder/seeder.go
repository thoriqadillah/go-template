package seeder

import (
	"context"
	"database/sql"
	"log"

	"github.com/uptrace/bun"
)

type Seeder interface {
	Seed(ctx context.Context, db *bun.Tx) error
}

// INFO: register all the seeder here
var seeders = make([]Seeder, 0)

func register(s ...Seeder) {
	seeders = append(seeders, s...)
}

func Seed(ctx context.Context, db *bun.DB) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Println("Error opening tx: ", err)
		return err
	}

	for _, seeder := range seeders {
		if err := seeder.Seed(ctx, &tx); err != nil {
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
