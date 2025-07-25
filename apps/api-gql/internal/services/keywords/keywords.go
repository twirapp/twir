package keywords

import (
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	KeywordsRepository keywords.Repository
	Logger             logger.Logger
	KeywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
}

func New(opts Opts) *Service {
	return &Service{
		keywordsRepository: opts.KeywordsRepository,
		logger:             opts.Logger,
		keywordsCacher:     opts.KeywordsCacher,
	}
}

const MaxPerChannel = 25

type Service struct {
	keywordsRepository keywords.Repository
	logger             logger.Logger
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
	}
}
