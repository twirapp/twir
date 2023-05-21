package types

import (
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
)

type ParseContextSender struct {
	ID          string
	Name        string
	DisplayName string
	Badges      []string
}

type ParseContextChannel struct {
	ID   string
	Name string
}

type ParseContext struct {
	Channel *ParseContextChannel
	Sender  *ParseContextSender

	Text      *string
	IsCommand bool

	Services *services.Services

	Cacher DataCacher
}
