package audit

import (
	"context"
)

// Recorder records audit information about operations that occurs in different applications.
type Recorder interface {
	RecordCreateOperation(ctx context.Context, operation CreateOperation) error
	RecordDeleteOperation(ctx context.Context, operation DeleteOperation) error
	RecordUpdateOperation(ctx context.Context, operation UpdateOperation) error
}
