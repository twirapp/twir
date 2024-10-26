package audit

import (
	"context"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/guregu/null"
	slogcommon "github.com/samber/slog-common"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/levels"
	"gorm.io/gorm"
)

func NewGorm(gormClient *gorm.DB) *AuditGorm {
	return &AuditGorm{
		gormClient: gormClient,
	}
}

func NewGormFx() func(gormClient *gorm.DB) *AuditGorm {
	return func(gormClient *gorm.DB) *AuditGorm {
		return NewGorm(gormClient)
	}
}

var _ slog.Handler = (*AuditGorm)(nil)

type AuditGorm struct {
	gormClient *gorm.DB

	attrs  []slog.Attr
	groups []string
}

// Enabled implements slog.Handler.
func (a *AuditGorm) Enabled(ctx context.Context, level slog.Level) bool {
	if level == levels.LevelAudit {
		return true
	}
	return false
}

// Handle implements slog.Handler.
func (a *AuditGorm) Handle(ctx context.Context, _ slog.Record) error {
	data, ok := ctx.Value(AuditFieldsContextKey{}).(Fields)
	if !ok {
		return nil
	}

	var oldValue null.String
	if data.OldValue != nil {
		oldValueBytes, _ := json.Marshal(data.OldValue)
		oldValue = null.StringFrom(string(oldValueBytes))
	}

	var newValue null.String
	if data.NewValue != nil {
		newValueBytes, _ := json.Marshal(data.NewValue)
		newValue = null.StringFrom(string(newValueBytes))
	}

	operationType := model.AuditOperationType(data.OperationType)
	if operationType == "" {
		operationType = model.AuditOperationTypeUnknown
	}

	a.gormClient.Save(
		&model.AuditLog{
			Table:         data.System,
			OperationType: operationType,
			OldValue:      oldValue,
			NewValue:      newValue,
			ObjectID:      null.StringFromPtr(data.ObjectID),
			ChannelID:     null.StringFromPtr(data.ChannelID),
			UserID:        null.StringFromPtr(data.ActorID),
		},
	)

	return nil
}

// WithAttrs implements slog.Handler.
func (a *AuditGorm) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AuditGorm{
		gormClient: a.gormClient,
		attrs:      slogcommon.AppendAttrsToGroup(a.groups, a.attrs, attrs...),
	}
}

// WithGroup implements slog.Handler.
func (a *AuditGorm) WithGroup(name string) slog.Handler {
	if name == "" {
		return a
	}

	return &AuditGorm{
		gormClient: a.gormClient,
		attrs:      a.attrs,
		groups:     append(a.groups, name),
	}
}
