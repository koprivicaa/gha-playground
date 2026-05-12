# Exercise 2 — Worker pool with error aggregation

## Problem

Implement a worker pool with this signature:

```go
// RunTasks runs every task with at most `concurrency` workers in flight.
// It collects every error and returns the slice. Order is not important.
// If ctx is cancelled, in-flight workers should return as soon as they notice.
func RunTasks(ctx context.Context, concurrency int, tasks []Task) []error
```

Where:

```go
type Task func(ctx context.Context) error
```

## Requirements

1. Use **buffered channel** as a job queue (or unbuffered with N goroutines pulling — both work).
2. Use `sync.WaitGroup` to wait for completion.
3. Errors must be collected into a slice that is **safe to write from multiple goroutines** (mutex or a dedicated channel).
4. If `ctx` is cancelled, no new tasks should start; running tasks should be told via the ctx they receive.
5. If `concurrency <= 0`, fall back to `concurrency = 1`.

## Stretch goals

- Add a fail-fast variant that cancels the context on the first error.
- Add a retry wrapper: `WithRetry(task Task, attempts int, backoff time.Duration) Task`.

## Acceptance test

A skeleton test is in `main_test.go`. Make it pass:

```sh
go test -race -v ./...
```

## Hints

- Pattern: spawn `concurrency` goroutines that range over a `chan Task`. After enqueuing all tasks, close the channel. Wait on the WaitGroup. Then return the errors.
- For thread-safety: a `sync.Mutex` around `[]error` is simplest. A second errors channel that you drain at the end is also clean.
- Make sure to test with `-race` to catch data races.
- For ctx cancellation: select on `<-ctx.Done()` inside the worker loop alongside the job channel.

## Talking points

- "I always reach for `errgroup.WithContext` from `golang.org/x/sync/errgroup` in production — it does this and gives you fail-fast for free."
- "Worker pool vs. unbounded goroutines: unbounded goroutines are fine when work is bounded and predictable, but for IO-heavy fan-out you want a pool to avoid resource exhaustion."
- "The `-race` flag in `go test` is non-negotiable for concurrency code."
