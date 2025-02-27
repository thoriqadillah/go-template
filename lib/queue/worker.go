package queue

import (
	"app/lib/notifier"

	"github.com/riverqueue/river"
)

type Worker interface {
	CreateWorker(workers *river.Workers)
}

var workers = river.NewWorkers()

func init() {
	processors := []Worker{
		notifier.CreateEmailWorker(),
	}

	for _, processor := range processors {
		processor.CreateWorker(workers)
	}
}
