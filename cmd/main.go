package main

import (
	"app/cmd/command"
	"app/db"
	"app/db/seeder"
	"app/env"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGABRT, syscall.SIGTERM)
	defer stop()

	godotenv.Load()

	bobdb, pool, err := db.Connect(ctx, env.DB_URL)
	if err != nil {
		panic(err)
	}

	defer bobdb.Close()
	defer pool.Close()

	app := &command.App{
		BobDB: bobdb,
	}

	// INFO: register all the command here
	command.Register(
		seeder.CreateCommand,
	)

	if err := command.Execute(ctx, app); err != nil {
		log.Fatal(err)
	}
}
