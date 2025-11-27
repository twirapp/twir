package recorder

import (
	"context"

	"github.com/twirapp/twir/libs/audit"
	"golang.org/x/sync/errgroup"
)

// Fanout is an [audit.Recorder] implementation that fans-out operations to multiple recorders in separate
// error-grouped goroutines.
type Fanout struct {
	recorders []audit.Recorder
}

var _ audit.Recorder = (*Fanout)(nil)

func NewFanout(recorders ...audit.Recorder) Fanout {
	return Fanout{
		recorders: recorders,
	}
}

func (f Fanout) RecordCreateOperation(ctx context.Context, operation audit.CreateOperation) error {
	return f.fanout(
		ctx, func(r audit.Recorder) error {
			return r.RecordCreateOperation(ctx, operation)
		},
	)
}

func (f Fanout) RecordDeleteOperation(ctx context.Context, operation audit.DeleteOperation) error {
	return f.fanout(
		ctx, func(r audit.Recorder) error {
			return r.RecordDeleteOperation(ctx, operation)
		},
	)
}

func (f Fanout) RecordUpdateOperation(ctx context.Context, operation audit.UpdateOperation) error {
	return f.fanout(
		ctx, func(r audit.Recorder) error {
			return r.RecordUpdateOperation(ctx, operation)
		},
	)
}

func (f Fanout) fanout(ctx context.Context, record func(recorder audit.Recorder) error) error {
	group, _ := errgroup.WithContext(ctx)

	for _, recorder := range f.recorders {
		group.Go(
			func() error {
				return record(recorder)
			},
		)
	}

	return group.Wait()
}
