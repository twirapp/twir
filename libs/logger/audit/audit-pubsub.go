package audit

import (
	"context"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/guregu/null"
	slogcommon "github.com/samber/slog-common"
	"github.com/twirapp/twir/libs/logger/levels"
	auditlogs "github.com/twirapp/twir/libs/pubsub/audit-logs"
	buscoreauditlogs "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

func NewPubSub(pb auditlogs.PubSub) *AuditPubsub {
	return &AuditPubsub{pb: pb}
}

func NewPubsubFx() func(pb auditlogs.PubSub) *AuditPubsub {
	return func(l auditlogs.PubSub) *AuditPubsub {
		return NewPubSub(l)
	}
}

var _ slog.Handler = (*AuditPubsub)(nil)

type AuditPubsub struct {
	pb auditlogs.PubSub

	attrs  []slog.Attr
	groups []string
}

func (c *AuditPubsub) Enabled(ctx context.Context, level slog.Level) bool {
	if level == levels.LevelAudit {
		return true
	}
	return false
}

func (c *AuditPubsub) Handle(ctx context.Context, record slog.Record) error {
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

	c.pb.Publish(
		ctx,
		auditlogs.AuditLog{
			ID:            uuid.New(),
			Table:         data.System,
			OperationType: buscoreauditlogs.AuditOperationType(data.OperationType),
			OldValue:      oldValue,
			NewValue:      newValue,
			ObjectID:      null.StringFromPtr(data.ObjectID),
			ChannelID:     null.StringFromPtr(data.ChannelID),
			UserID:        null.StringFromPtr(data.ActorID),
			CreatedAt:     time.Now(),
		},
	)

	return nil
}

func (c *AuditPubsub) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AuditPubsub{
		pb:    c.pb,
		attrs: slogcommon.AppendAttrsToGroup(c.groups, c.attrs, attrs...),
	}
}

func (c *AuditPubsub) WithGroup(name string) slog.Handler {
	if name == "" {
		return c
	}

	return &AuditPubsub{
		pb:     c.pb,
		attrs:  c.attrs,
		groups: append(c.groups, name),
	}
}
