package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-pg/pg"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/gitbook"
)

func Init(db *pg.DB) func(context.Context, *sync.WaitGroup) {
	return func(ctx context.Context, waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()

		scheduler, errScheduler := gocron.NewScheduler()
		if errScheduler != nil {
			panic(errScheduler)
		}

		_, errorJob := scheduler.NewJob(
			gocron.DurationJob(
				15*time.Second,
			),
			gocron.NewTask(
				func() {
					err := gitbook.Vectorize(db)
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
