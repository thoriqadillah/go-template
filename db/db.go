package db

import (
	"app/lib/log"
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stephenafamo/bob"
)

var sqldb *sql.DB
var db *bob.DB
var logger = log.Logger()

func Db() *bob.DB {
	return db
}

func Sql() *sql.DB {
	return sqldb
}

func Connect(connstr string) (closer func(), err error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		logger.Fatal(err.Error())
		return nil, err
	}

	sqldb = stdlib.OpenDBFromPool(pool)
	bob := bob.NewDB(sqldb)

	db = &bob
	close := func() {
		pool.Close()
		db.Close()
	}

	return close, sqldb.Ping()
}
