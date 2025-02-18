package db

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
		return
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
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
		log.Fatalf("Could not start resource: %s", err)
	}

	pg.Expire(30)
	hostPort := pg.GetHostPort("5432/tcp")
	pgUrl := fmt.Sprintf("postgresql://postgres@%s/apptest?sslmode=disable", hostPort)

	// Wait for the Postgres to be ready
	if err := pool.Retry(func() error {
		Connect(pgUrl)
		return sqldb.Ping()
	}); err != nil {
		panic("Could not connect to postgres: " + err.Error())
	}

	defer func() {
		if err := pool.Purge(pg); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestMigration(t *testing.T) {
	err := goose.SetDialect("postgres")
	assert.NoError(t, err)

	pwd, _ := os.Getwd()

	err = goose.Up(sqldb, pwd+"/migration")
	assert.NoError(t, err)
}
