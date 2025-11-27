package keywords

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	KeywordsRepository keywords.Repository
	AuditRecorder      audit.Recorder
	Logger             *slog.Logger
	KeywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
}

func New(opts Opts) *Service {
	return &Service{
		keywordsRepository: opts.KeywordsRepository,
		auditRecorder:      opts.AuditRecorder,
		logger:             opts.Logger,
		keywordsCacher:     opts.KeywordsCacher,
	}
}

const MaxPerChannel = 25

type Service struct {
	keywordsRepository keywords.Repository
	auditRecorder      audit.Recorder
	logger             *slog.Logger
	keywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
}

func (c *Service) dbToModel(m model.Keyword) entity.Keyword {
	return entity.Keyword{
		ID:               m.ID,
		ChannelID:        m.ChannelID,
		Text:             m.Text,
		Response:         m.Response,
		Enabled:          m.Enabled,
		Cooldown:         m.Cooldown,
		CooldownExpireAt: m.CooldownExpireAt,
		IsReply:          m.IsReply,
		IsRegular:        m.IsRegular,
		Usages:           m.Usages,
		RolesIDs:         m.RolesIDs,
	}
}
