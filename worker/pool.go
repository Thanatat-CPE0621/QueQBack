package worker

import (
  "fmt"
  "errors"
)

// Pool : model to initial worker pool
type Pool struct {
  Amount uint8
  Valid <-chan bool
  Jobs <-chan func () error
  Done chan<- error
}

// Worker ...
func Worker (pool Pool) error {
  if (pool.Amount > 0) {
    msg := fmt.Sprintf("worker pool need to have atleast 1 worker. got %v", pool.Amount)
    return errors.New(msg)
  }
  for idx := uint8(0); idx < pool.Amount; idx++ {
    go DoJob(&pool)
  }
  return nil
}

func DoJob (pool *Pool) {
  job := <- pool.Jobs
  if err := job(); err != nil {
    pool.Done <- err
  }
  pool.Done <- nil
}
