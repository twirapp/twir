package uploadedfiles

import (
	"context"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
)

var ErrNotFound = errors.New("uploaded file not found")

type CreateInput struct {
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
}

type GetListInput struct {
	UserID  string
	Page    int
	PerPage int
}

type GetListOutput struct {
	Items []uploadedfile.Entity
	Total int
}

type Repository interface {
	Create(ctx context.Context, input CreateInput) (uploadedfile.Entity, error)
	GetByPublicID(ctx context.Context, publicID string) (uploadedfile.Entity, error)
	GetManyByPublicIDs(ctx context.Context, publicIDs []string) ([]uploadedfile.Entity, error)
	GetList(ctx context.Context, input GetListInput) (GetListOutput, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetExpired(ctx context.Context, limit int) ([]uploadedfile.Entity, error)
}
