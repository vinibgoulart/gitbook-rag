package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-pg/pg"
	_ "github.com/joho/godotenv/autoload"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/database"
	synchronizer "github.com/vinibgoulart/gitbook-postgresql-vectorize/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "gitbook_postgresql_vectorize",
	})
	defer db.Close()

	err := database.CreateSchemaDatabase(db)
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go synchronizer.Init(ctx, &waitGroup)

	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)

	<-closeChannel
	cancel()

	waitGroup.Wait()
}
