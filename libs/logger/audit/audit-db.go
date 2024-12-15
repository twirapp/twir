package audit

import (
	"context"
	"log/slog"

	"github.com/goccy/go-json"
	slogcommon "github.com/samber/slog-common"
	"github.com/satont/twir/libs/logger/levels"
	audit_logs "github.com/twirapp/twir/libs/repositories/audit-logs"
	"github.com/twirapp/twir/libs/repositories/audit-logs/model"
)

func NewDatabase(repository audit_logs.Repository) *AuditDatabase {
	return &AuditDatabase{
		repository: repository,
	}
}

func NewDatabaseFx(repository audit_logs.Repository) *AuditDatabase {
	return NewDatabase(repository)
}

var _ slog.Handler = (*AuditDatabase)(nil)

type AuditDatabase struct {
	repository audit_logs.Repository

	attrs  []slog.Attr
	groups []string
}

// Enabled implements slog.Handler.
func (a *AuditDatabase) Enabled(ctx context.Context, level slog.Level) bool {
	if level == levels.LevelAudit {
		return true
	}
	return false
}

// Handle implements slog.Handler.
func (a *AuditDatabase) Handle(ctx context.Context, _ slog.Record) error {
	data, ok := ctx.Value(AuditFieldsContextKey{}).(Fields)
	if !ok {
		return nil
	}

	var oldValue *string
	if data.OldValue != nil {
		oldValueBytes, _ := json.Marshal(data.OldValue)
		oldValueStr := string(oldValueBytes)
		oldValue = &oldValueStr
	}

	var newValue *string
	if data.NewValue != nil {
		newValueBytes, _ := json.Marshal(data.NewValue)
		newValueStr := string(newValueBytes)
		newValue = &newValueStr
	}

	operationType := model.AuditOperationType(data.OperationType)
	if operationType == "" {
		operationType = model.AuditOperationTypeUnknown
	}

	a.repository.Create(
		ctx, audit_logs.CreateInput{
			Table:         data.System,
			OperationType: operationType,
			OldValue:      oldValue,
			NewValue:      newValue,
			ObjectID:      data.ObjectID,
			ChannelID:     data.ChannelID,
			UserID:        data.ActorID,
		},
	)

	return nil
}

// WithAttrs implements slog.Handler.
func (a *AuditDatabase) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AuditDatabase{
		repository: a.repository,
		attrs:      slogcommon.AppendAttrsToGroup(a.groups, a.attrs, attrs...),
	}
}

// WithGroup implements slog.Handler.
func (a *AuditDatabase) WithGroup(name string) slog.Handler {
	if name == "" {
		return a
	}

	return &AuditDatabase{
		repository: a.repository,
		attrs:      a.attrs,
		groups:     append(a.groups, name),
	}
}
