package db

import (
	"app/env"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	_ "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var sqldb *sql.DB
var bundb *bun.DB

func Db() *bun.DB {
	return bundb
}

func Connect(ctx context.Context, connstr string) (*bun.DB, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, nil, err
	}

	sqldb = stdlib.OpenDBFromPool(pool)
	bundb = bun.NewDB(sqldb, pgdialect.New())
	bundb.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(env.DEV),
	))

	return bundb, pool, sqldb.Ping()
}

func Migrate(pathDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(sqldb, pathDir)
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
			"POSTGRES_DB=packformtest",
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
	pgUrl := fmt.Sprintf("postgresql://postgres@%s/packformtest?sslmode=disable", hostPort)

	// Wait for the Postgres to be ready
	err = pool.Retry(func() error {
		_, _, err := Connect(ctx, pgUrl)
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
