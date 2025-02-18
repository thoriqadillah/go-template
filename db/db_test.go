package db

import (
	"os"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	purge := SetupTest()
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
