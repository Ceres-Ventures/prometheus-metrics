package job

import "fmt"

type (
	Dispatcher struct {
		jobCounter     uint64
		jobQueue       chan *Job
		dispatchStatus chan *DispatchStatus
		workQueue      chan *Job
		workerQueue    chan *Worker
	}

	DispatchStatus struct {
		Type   string
		ID     uint64
		Status string
	}
)

func CreateNewDispatcher() *Dispatcher {
	return &Dispatcher{
		jobCounter:     0,
		jobQueue:       make(chan *Job),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue:      make(chan *Job),
		workerQueue:    make(chan *Worker),
	}
}

func (d *Dispatcher) Start(numWorkers uint64) {
	for i := uint64(0); i < numWorkers; i++ {
		worker := CreateNewWorker(d, i, d.workerQueue, d.workQueue, d.dispatchStatus)
		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-d.jobQueue:
				d.workQueue <- job
			case ds := <-d.dispatchStatus:
				fmt.Printf("Got a dispatch status:\n\tType[%s] - ID[%d] - Status[%s]\n", ds.Type, ds.ID, ds.Status)
				if ds.Type == "worker" {
					if ds.Status == "quit" {
						d.jobCounter--
					}
				}
			}
		}
	}()
}

func (d *Dispatcher) AddJob(je JobExecutable, repeat bool, repeatCount int, repeatCounter int) {
	j := &Job{
		ID:            d.jobCounter,
		F:             je,
		RepeatCount:   repeatCount,
		RepeatCounter: repeatCounter,
		IsRepeatable:  repeat,
	}
	go func() {
		d.jobQueue <- j
	}()
	d.jobCounter++
}

func (d *Dispatcher) Finished() bool {
	if d.jobCounter < 1 {
		return true
	}
	return false
}
