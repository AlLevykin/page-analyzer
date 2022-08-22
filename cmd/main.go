package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"page-analyzer/internal/analyzer"
	"page-analyzer/internal/urls"
	"page-analyzer/internal/wpool"
	"syscall"
)

const workerCount = 4

func main() {

	fmt.Println("processing...")

	if len(os.Args) != 2 {
		log.Fatal("wrong arguments")
	}

	us, err := urls.Parse(os.Args[1], ",")
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}

	jobsCount := len(us)
	jobs := make([]wpool.Job, jobsCount)
	for i := 0; i < jobsCount; i++ {
		jobs[i] = wpool.Job{
			Descriptor: wpool.JobDescriptor{
				ID: wpool.JobID(us[i]),
			},
			ExecFn: analyzer.Process,
			Args:   us[i],
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	wp := wpool.NewWorkerPool(workerCount)
	go wp.GenerateFrom(jobs)
	go wp.Run(ctx)

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				return
			}
			ResDispatcher(r)
		}
	}

}

func ResDispatcher(r wpool.Result) {
	if r.Err != nil {
		fmt.Printf("%s: %v\n", r.Descriptor.ID, r.Err)
		return
	}

	val, ok := r.Value.(string)
	if !ok {
		fmt.Println("wrong result type")
		return
	}
	fmt.Println(val)
}
