package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksCh <- task
	}
	errCh := make(chan error, len(tasks))
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case task := <-tasksCh:
					if len(errCh) > m {
						return
					}
					err := task()
					if err != nil {
						errCh <- err
					}
				default:
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)
	if len(errCh) >= m {
		// case, when m == 0
		if err := <-errCh; err == nil {
			return nil
		}
		return ErrErrorsLimitExceeded
	}

	return nil
}
