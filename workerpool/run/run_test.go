package run

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunTasks_AllSucceed(t *testing.T) {
	var counter int64
	tasks := make([]Task, 100)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			atomic.AddInt64(&counter, 1)
			return nil
		}
	}
	errs := RunTasks(context.Background(), 10, tasks)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if got := atomic.LoadInt64(&counter); got != 100 {
		t.Fatalf("expected 100 runs, got %d", got)
	}
}

func TestRunTasks_CollectsErrors(t *testing.T) {
	tasks := []Task{
		func(ctx context.Context) error { return errors.New("a") },
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return errors.New("c") },
	}
	errs := RunTasks(context.Background(), 2, tasks)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(errs))
	}
}

func TestRunTasks_Cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled
	called := int64(0)
	tasks := make([]Task, 50)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			atomic.AddInt64(&called, 1)
			time.Sleep(10 * time.Millisecond)
			return nil
		}
	}
	_ = RunTasks(ctx, 4, tasks)
	if atomic.LoadInt64(&called) == int64(len(tasks)) {
		t.Fatalf("expected ctx cancellation to short-circuit some tasks; all %d ran", len(tasks))
	}
}

func TestRunTasks_ZeroConcurrencyFallsBackToOne(t *testing.T) {
	var counter int64
	tasks := make([]Task, 5)
	for i := range tasks {
		tasks[i] = func(ctx context.Context) error {
			atomic.AddInt64(&counter, 1)
			return nil
		}
	}
	errs := RunTasks(context.Background(), 0, tasks)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if got := atomic.LoadInt64(&counter); got != 5 {
		t.Fatalf("expected 5 runs, got %d", got)
	}
}

func TestRunTasksFailFast_StopsOnFirstError(t *testing.T) {
	var ran int64
	// First task fails immediately; the rest just count and sleep so
	// cancellation has time to reach them.
	tasks := []Task{
		func(ctx context.Context) error { return errors.New("boom") },
	}
	for i := 0; i < 20; i++ {
		tasks = append(tasks, func(ctx context.Context) error {
			atomic.AddInt64(&ran, 1)
			time.Sleep(20 * time.Millisecond)
			return nil
		})
	}

	errs := RunTasksFailFast(context.Background(), 2, tasks)
	if len(errs) == 0 {
		t.Fatal("expected at least one error")
	}
	if got := atomic.LoadInt64(&ran); got == 20 {
		t.Fatalf("expected fail-fast to stop some tasks, but all 20 ran")
	}
}

func TestWithRetry_SucceedsAfterRetries(t *testing.T) {
	var calls int64
	// Fails twice, succeeds on the third attempt.
	task := WithRetry(func(ctx context.Context) error {
		n := atomic.AddInt64(&calls, 1)
		if n < 3 {
			return errors.New("not yet")
		}
		return nil
	}, 3, time.Millisecond)

	if err := task(context.Background()); err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if got := atomic.LoadInt64(&calls); got != 3 {
		t.Fatalf("expected 3 calls, got %d", got)
	}
}

func TestWithRetry_ReturnsLastError(t *testing.T) {
	task := WithRetry(func(ctx context.Context) error {
		return errors.New("always fails")
	}, 3, time.Millisecond)

	if err := task(context.Background()); err == nil {
		t.Fatal("expected error after exhausting retries")
	}
}

func TestWithRetry_RespectsCtxCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	var calls int64
	task := WithRetry(func(ctx context.Context) error {
		atomic.AddInt64(&calls, 1)
		return errors.New("fail")
	}, 5, 50*time.Millisecond)

	// With a cancelled ctx the backoff select should return ctx.Err()
	// immediately, so the task should not be retried 5 times.
	task(ctx)
	if got := atomic.LoadInt64(&calls); got == 5 {
		t.Fatalf("expected retry to abort early on cancelled ctx, but ran %d times", got)
	}
}
