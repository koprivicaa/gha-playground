package run

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Execute(ctx context.Context) error {
	// Run tasks with concurrency of 5 and print any errors.
	tasks := []Task{
		func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return fmt.Errorf("task 1 failed")
		},
		func(ctx context.Context) error {
			time.Sleep(50 * time.Millisecond)
			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(150 * time.Millisecond)
			return fmt.Errorf("task 3 failed")
		},
	}
	errs := RunTasks(ctx, 5, tasks)
	
	// collect errors into a single error for simplicity
	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors: %v", len(errs), errs)
	}
	return nil
}

// Task is a unit of work.
type Task func(ctx context.Context) error

// RunTasks runs the tasks with at most `concurrency` workers and aggregates errors.
func RunTasks(ctx context.Context, concurrency int, tasks []Task) []error {
	// HINT 1: handle concurrency <= 0
	if concurrency <= 0 {
		concurrency = 1
	}
	// HINT 2: tasksChan := make(chan Task, len(tasks))
	tasksChan := make(chan Task)
	errorsChan := make(chan error)

	//Send tasks
	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			tasksChan <- task
		}
	}()

	var wg sync.WaitGroup
	// HINT 3: spawn N goroutines that select on <-tasksChan and <-ctx.Done()
	//Start workers UNBAFFERED CHANNEL
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case task, ok := <-tasksChan:
						if !ok {
							return
						}
						if err := task(ctx); err != nil {
							errorsChan <- err
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}

	//Close errorsChan after all workers are done
	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	//Collect errors
	var (
		errs []error
		mu   sync.Mutex
	)
	
	for err := range errorsChan {
		if err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}
		
	return errs
}

// RunTasksFailFast is like RunTasks but cancels all remaining work on the first error.
func RunTasksFailFast(ctx context.Context, concurrency int, tasks []Task) []error {
	// Derive a child context so we can cancel it on first error
	// without affecting the caller's context.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if concurrency <= 0 {
		concurrency = 1
	}

	tasksChan := make(chan Task, len(tasks))
	errorsChan := make(chan error, len(tasks))

	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			tasksChan <- task
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasksChan:
					if !ok {
						return
					}
					if err := task(ctx); err != nil {
						errorsChan <- err
						// Cancel the child context so other workers stop
						// picking up new tasks after this point.
						cancel()
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errorsChan)
	}()

	var errs []error
	for err := range errorsChan {
		errs = append(errs, err)
	}
	return errs
}

// WithRetry wraps a task so it is retried up to `attempts` times.
// It waits `backoff` duration between attempts and respects ctx cancellation.
func WithRetry(task Task, attempts int, backoff time.Duration) Task {
	return func(ctx context.Context) error {
		var err error
		for i := 0; i < attempts; i++ {
			err = task(ctx)
			if err == nil {
				return nil
			}
			// Don't wait after the last attempt.
			if i < attempts-1 {
				select {
				case <-time.After(backoff):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return err
	}
}