package synchronizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/vinibgoulart/gitbook-postgresql-vectorize/packages/gitbook"
)

func Init(ctx context.Context, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	scheduler, errScheduler := gocron.NewScheduler()
	if errScheduler != nil {
		panic(errScheduler)
	}

	_, errorJob := scheduler.NewJob(
		gocron.DurationJob(
			4*time.Second,
		),
		gocron.NewTask(
			func() {
				err := gitbook.Vectorize()
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
