package bot

import (
	"context"
)

type Bot interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	SendMessage(ctx context.Context, input SendMessageInput) error
	SendPhoto(ctx context.Context, input SendPhotoInput) error
}
