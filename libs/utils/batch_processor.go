package utils

import (
	"context"
	"sync"
	"time"
)

type BatchProcessorFn[T any] func(ctx context.Context, data []T)

type BatchProcessor[T any] struct {
	callback BatchProcessorFn[T]
	queue    chan T

	shutdown   chan struct{}
	shutdownWg sync.WaitGroup

	batchSize int
	interval  time.Duration
}

type BatchProcessorOpts[T any] struct {
	Interval  time.Duration
	BatchSize int
	Callback  BatchProcessorFn[T]
}

func NewBatchProcessor[T any](opts BatchProcessorOpts[T]) *BatchProcessor[T] {
	if opts.BatchSize <= 0 {
		opts.BatchSize = 100
	}

	return &BatchProcessor[T]{
		callback:  opts.Callback,
		queue:     make(chan T, opts.BatchSize),
		shutdown:  make(chan struct{}),
		batchSize: opts.BatchSize,
		interval:  opts.Interval,
	}
}

func (s *BatchProcessor[T]) Add(data T) bool {
	select {
	case <-s.shutdown:
		return false
	case s.queue <- data:
		return true
	}
}

func (s *BatchProcessor[T]) Start(ctx context.Context) {
	s.shutdownWg.Add(1)
	defer s.shutdownWg.Done()

	var (
		batch = make([]T, 0, s.batchSize)
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.shutdown:
			s.processBatchQueue(ctx, batch)
			return
		default:
			time.Sleep(s.interval)

			s.processBatchQueue(ctx, batch)
			batch = batch[:0]
		}
	}
}

func (s *BatchProcessor[T]) Shutdown(ctx context.Context) error {
	select {
	case <-s.shutdown:
		return nil
	default:
		s.shutdown <- struct{}{}
		close(s.shutdown)
	}

	done := make(chan struct{})

	go func() {
		s.shutdownWg.Wait()
		done <- struct{}{}
		close(done)
	}()

	select {
	case <-done:
		close(s.queue)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *BatchProcessor[T]) processBatchQueue(ctx context.Context, batch []T) {
	for len(s.queue) > 0 {
		batch = append(batch, <-s.queue)
		if len(batch) == s.batchSize {
			break
		}
	}

	if len(batch) == 0 {
		return
	}

	s.callback(ctx, batch)
}
