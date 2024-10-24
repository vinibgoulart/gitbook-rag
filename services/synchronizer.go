package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/uptrace/bun"
	"github.com/vinibgoulart/gitbook-rag/packages/gitbook"
)

func Init(db *bun.DB) func(context.Context, *sync.WaitGroup) {
	return func(ctx context.Context, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		fmt.Println("Synchronizer started")

		scheduler, errScheduler := gocron.NewScheduler()
		if errScheduler != nil {
			panic(errScheduler)
		}

		_, errorJob := scheduler.NewJob(
			gocron.DurationJob(
				1*time.Hour,
			),
			gocron.NewTask(
				func() {
					err := gitbook.Vectorize(&ctx, db)
					if err != nil {
						fmt.Println(err.Error())
					}
				},
			),
		)

		if errorJob != nil {
			panic(errorJob)
		}

		scheduler.Start()

		<-ctx.Done()

		_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
	}
}
