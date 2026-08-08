package uploadedfile

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID               uuid.UUID
	PublicID         string
	UploadedByUserID *string
	FileName         *string
	MimeType         string
	Extension        string
	SizeBytes        int64
	S3Key            string
	DeleteKey        string
	UserAgent        *string
	UserIP           *netip.Addr
	ExpiresAt        time.Time
	CreatedAt        time.Time

	isNil bool
}

func (e Entity) IsNil() bool {
	return e.isNil
}

var Nil = Entity{isNil: true}
