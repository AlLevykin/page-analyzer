package wpool

import (
	"context"
	"errors"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			return
		}
	}
}

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func NewWorkerPool(wcount int) WorkerPool {
	return WorkerPool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}

func (wp WorkerPool) Run(ctx context.Context) {
	defer func() {
		close(wp.Done)
		close(wp.results)
	}()

	var wg sync.WaitGroup
	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}
	wg.Wait()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			wp.results <- Result{
				Err:        errors.New("cancel"),
				Descriptor: job.Descriptor,
			}
		}
	}
}

func (wp WorkerPool) Results() <-chan Result {
	return wp.results
}

func (wp WorkerPool) GenerateFrom(jobsBulk []Job) {
	for i, _ := range jobsBulk {
		wp.jobs <- jobsBulk[i]
	}
	close(wp.jobs)
}
