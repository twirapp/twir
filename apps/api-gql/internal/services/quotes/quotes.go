package quotes

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/quotes"
	"github.com/twirapp/twir/libs/repositories/quotes/model"
)

func New(
	quotesRepository quotes.Repository,
	auditRecorder audit.Recorder,
	logger *slog.Logger,
	quotesCacher *generic_cacher.GenericCacher[[]model.Quote],
) *Service {
	return &Service{
		quotesRepository: quotesRepository,
		auditRecorder:    auditRecorder,
		logger:           logger,
		quotesCacher:     quotesCacher,
	}
}

type Service struct {
	quotesRepository quotes.Repository
	auditRecorder    audit.Recorder
	logger           *slog.Logger
	quotesCacher     *generic_cacher.GenericCacher[[]model.Quote]
}

func (c *Service) dbToModel(m model.Quote) entity.Quote {
	return entity.Quote{
		ID:          m.ID,
		ChannelID:   m.ChannelID.String(),
		Number:      m.Number,
		Text:        m.Text,
		CreatorID:   m.CreatorID,
		CreatorName: m.CreatorName,
		GameID:      m.GameID,
		GameName:    m.GameName,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
