package quotes

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/quotes"
	"github.com/twirapp/twir/libs/repositories/quotes/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	QuotesRepository quotes.Repository
	AuditRecorder    audit.Recorder
	Logger           *slog.Logger
	QuotesCacher     *generic_cacher.GenericCacher[[]model.Quote]
}

func New(opts Opts) *Service {
	return &Service{
		quotesRepository: opts.QuotesRepository,
		auditRecorder:    opts.AuditRecorder,
		logger:           opts.Logger,
		quotesCacher:     opts.QuotesCacher,
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
