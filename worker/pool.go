package worker

import (
	"errors"
	"fmt"
)

// Pool : model to initial worker pool
type Pool struct {
	Amount uint8
	Jobs   <-chan func()
	Done   chan<- error
}

// Worker ...
func Worker(pool Pool) error {
	if pool.Amount < 1 {
		msg := fmt.Sprintf("worker pool need to have atleast 1 worker. got %v", pool.Amount)
		return errors.New(msg)
	}
	for idx := uint8(0); idx < pool.Amount; idx++ {
		go doJob(&pool)
	}
	return nil
}

func doJob(pool *Pool) {
	for job := range pool.Jobs {
		job()
	}
}
