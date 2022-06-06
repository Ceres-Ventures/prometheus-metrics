package job

import (
	"fmt"

	"github.com/labstack/gommon/log"
)

type (
	Worker struct {
		ID             uint64
		jobs           chan *Job
		dispatchStatus chan *DispatchStatus
		Quit           chan bool
		Dispatcher     *Dispatcher
	}
)

func CreateNewWorker(d *Dispatcher, id uint64, workerQueue chan *Worker, jobQueue chan *Job, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		Dispatcher:     d,
		ID:             id,
		jobs:           jobQueue,
		dispatchStatus: dStatus,
	}

	go func() {
		workerQueue <- w
	}()
	return w
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.jobs:
				err := job.F()
				if err != nil {
					log.Warn(fmt.Sprintf("Job has errored out: %s", err.Error()))
				}
				job.RepeatCounter++
				w.dispatchStatus <- &DispatchStatus{Type: "worker", ID: w.ID, Status: "quit"}
				// See if it is repeatable
				if job.IsRepeatable {
					if job.RepeatCounter < job.RepeatCount || job.RepeatCount == -1 {
						w.Dispatcher.AddJob(job.F, job.IsRepeatable, job.RepeatCount, job.RepeatCounter)
					} else {
						w.Quit <- true
					}
				} else {
					w.Quit <- true
				}
			case <-w.Quit:
				return
			}
		}
	}()
}
