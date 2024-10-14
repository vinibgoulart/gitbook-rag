package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	synchronizer "github.com/vinibgoulart/gitbook-postgresql-vectorize/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go synchronizer.Init(ctx, &waitGroup)

	closeChannel := make(chan os.Signal, 1)
	signal.Notify(closeChannel, syscall.SIGINT, syscall.SIGTERM)

	<-closeChannel
	cancel()

	waitGroup.Wait()
}
