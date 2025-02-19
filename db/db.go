package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stephenafamo/bob"
)

var sqldb *sql.DB
var db *bob.DB

func Db() *bob.DB {
	return db
}

func Connect(ctx context.Context, connstr string) (close func(), err error) {
	pool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	sqldb = stdlib.OpenDBFromPool(pool)
	bob := bob.NewDB(sqldb)

	db = &bob
	close = func() {
		pool.Close()
		db.Close()
	}

	return close, sqldb.Ping()
}

func SetupTest(ctx context.Context) (purge func()) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %v", err)
		return
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}

	pg, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17-alpine",
		Env: []string{
			"POSTGRES_HOST_AUTH_METHOD=trust",
			"POSTGRES_DB=apptest",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}

	pg.Expire(30)
	hostPort := pg.GetHostPort("5432/tcp")
	pgUrl := fmt.Sprintf("postgresql://postgres@%s/apptest?sslmode=disable", hostPort)

	// Wait for the Postgres to be ready
	err = pool.Retry(func() error {
		_, err := Connect(ctx, pgUrl)
		return err
	})

	if err != nil {
		log.Fatalf("Could not connect to postgres: %v", err)
	}

	return func() {
		if err := pool.Purge(pg); err != nil {
			log.Fatalf("Could not purge resource: %v", err)
		}
	}
}
