package keywords

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
	"github.com/twirapp/twir/libs/repositories/plans"
)

func New(
	keywordsRepository keywords.Repository,
	auditRecorder audit.Recorder,
	logger *slog.Logger,
	keywordsCacher *generic_cacher.GenericCacher[[]model.Keyword],
	plansRepository plans.Repository,
) *Service {
	return &Service{
		keywordsRepository: keywordsRepository,
		auditRecorder:      auditRecorder,
		logger:             logger,
		keywordsCacher:     keywordsCacher,
		plansRepository:    plansRepository,
	}
}

type Service struct {
	keywordsRepository keywords.Repository
	auditRecorder      audit.Recorder
	logger             *slog.Logger
	keywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
	plansRepository    plans.Repository
}

func (c *Service) dbToModel(m model.Keyword) entity.Keyword {
	return entity.Keyword{
		ID:               m.ID,
		ChannelID:        m.ChannelID.String(),
		Text:             m.Text,
		Response:         m.Response,
		Enabled:          m.Enabled,
		Cooldown:         m.Cooldown,
		CooldownExpireAt: m.CooldownExpireAt,
		IsReply:          m.IsReply,
		IsRegular:        m.IsRegular,
		Usages:           m.Usages,
		RolesIDs:         m.RolesIDs,
		Platforms:        m.Platforms,
	}
}
