package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workers int, maxErrors int) error {
	taskCh := make(chan Task, len(tasks))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskCh:
					if !ok {
						cancel()
						return
					}

					if task() != nil {
						mu.Lock()
						maxErrors--
						if maxErrors <= 0 {
							mu.Unlock()
							cancel()
							return
						}
						mu.Unlock()
					}
				default:
				}
			}
		}()
	}

	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if maxErrors <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
