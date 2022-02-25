package download

import (
	"context"
	"sync"
)

type Dispatcher struct {
	sem    chan struct{}
	queue  chan *Job
	worker Worker
	wg     sync.WaitGroup
}

func NewDispatcher(worker Worker, maxWorkers int, queueSize int) *Dispatcher {
	return &Dispatcher{
		sem:    make(chan struct{}, maxWorkers),
		queue:  make(chan *Job, queueSize),
		worker: worker,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	d.wg.Add(1)
	go d.loop(ctx)
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) AddJob(job *Job) {
	d.queue <- job
}

func (d *Dispatcher) AddJobs(jobs []*Job) {
	for _, job := range jobs {
		d.AddJob(job)
	}
}

func (d *Dispatcher) stop() {
	d.wg.Done()
}

func (d *Dispatcher) loop(ctx context.Context) {
	var wg sync.WaitGroup
Loop:
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			break Loop
		case job := <-d.queue:
			wg.Add(1)
			d.sem <- struct{}{}
			go func(job *Job) {
				defer wg.Done()
				defer func() { <-d.sem }()
				d.worker.RunJob(job)
			}(job)
		}
	}
	d.stop()
}
