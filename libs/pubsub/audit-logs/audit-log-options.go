package auditlog

import (
	"github.com/guregu/null"
)

type Option func(*AuditLog)

func WithOldValue(oldValue string) Option {
	return func(al *AuditLog) {
		al.OldValue = null.StringFrom(oldValue)
	}
}

func WithNewValue(newValue string) Option {
	return func(al *AuditLog) {
		al.NewValue = null.StringFrom(newValue)
	}
}

func WithObjectID(objectID string) Option {
	return func(al *AuditLog) {
		al.ObjectID = null.StringFrom(objectID)
	}
}

func WithChannelID(channelID string) Option {
	return func(al *AuditLog) {
		al.ChannelID = null.StringFrom(channelID)
	}
}

func WithUserID(userID string) Option {
	return func(al *AuditLog) {
		al.UserID = null.StringFrom(userID)
	}
}
