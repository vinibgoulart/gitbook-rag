package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/vinibgoulart/gitbook-rag/http"
	"github.com/vinibgoulart/gitbook-rag/packages/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(dsn)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	err := database.CreateSchemaDatabase(db, ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go http.ServerInit(db)(ctx, &waitGroup)

	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)

	<-closeChannel
	cancel()

	waitGroup.Wait()
}
