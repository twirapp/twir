package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/guregu/null"
)

type AuditLog struct {
	TableName     string
	OperationType AuditOperationType
	OldValue      null.String
	NewValue      null.String
	ObjectID      null.String
	ChannelID     null.String
	UserID        null.String
	CreatedAt     time.Time
}

type AuditOperationType string

func (a AuditOperationType) Value() (driver.Value, error) {
	return string(a), nil
}

var (
	errScanNilAuditOperationType = errors.New("Scan on nil *AuditOperationType")
)

func (a *AuditOperationType) Scan(src any) error {
	if a == nil {
		return errScanNilAuditOperationType
	}

	switch s := src.(type) {
	case string:
		*a = AuditOperationType(s)
	case []byte:
		*a = AuditOperationType(s)
	case nil:
		*a = ""
	default:
		return fmt.Errorf("unsupported Scan, storing %T into *AuditOperationType", src)
	}

	return nil
}

const (
	AuditOperationCreate      AuditOperationType = "CREATE"
	AuditOperationUpdate      AuditOperationType = "UPDATE"
	AuditOperationDelete      AuditOperationType = "DELETE"
	AuditOperationTypeUnknown AuditOperationType = "UNKNOWN"
)

var Nil = AuditLog{}
