package db

import (
	"context"
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestMain(m *testing.M) {
	purge := SetupTest(ctx)
	defer purge()

	m.Run()
}

func TestMigration(t *testing.T) {
	err := goose.SetDialect("postgres")
	assert.NoError(t, err)

	pwd, _ := os.Getwd()

	err = goose.Up(sqldb, pwd+"/migration")
	assert.NoError(t, err)
}
