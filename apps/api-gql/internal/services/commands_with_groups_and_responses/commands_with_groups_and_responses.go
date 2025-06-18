package commands_with_groups_and_responses

import (
	"context"
	"slices"
	"strings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	commandsservice "github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TrmManager                               trm.Manager
	CommandsRepository                       commands.Repository
	CommandsWithGroupsAndResponsesRepository commands_with_groups_and_responses.Repository
	ResponsesRepository                      commands_response.Repository
	CommandsService                          *commandsservice.Service
	Logger                                   logger.Logger
	CachedCommandsClient                     *generic_cacher.GenericCacher[[]deprecatedgormmodel.ChannelsCommands]
}

func New(opts Opts) *Service {
	return &Service{
		trmManager:                               opts.TrmManager,
		commandsWithGroupsAndResponsesRepository: opts.CommandsWithGroupsAndResponsesRepository,
		responsesRepository:                      opts.ResponsesRepository,
		commandsRepository:                       opts.CommandsRepository,
		logger:                                   opts.Logger,
		commandsService:                          opts.CommandsService,
		cachedCommandsClient:                     opts.CachedCommandsClient,
	}
}

type Service struct {
	trmManager                               trm.Manager
	commandsRepository                       commands.Repository
	commandsWithGroupsAndResponsesRepository commands_with_groups_and_responses.Repository
	responsesRepository                      commands_response.Repository

	commandsService      *commandsservice.Service
	logger               logger.Logger
	cachedCommandsClient *generic_cacher.GenericCacher[[]deprecatedgormmodel.ChannelsCommands]
}

func (c *Service) mapToEntity(m model.CommandWithGroupAndResponses) entity.CommandWithGroupAndResponses {
	e := entity.CommandWithGroupAndResponses{
		Command: entity.Command{
			ID:                        m.Command.ID,
			Name:                      m.Command.Name,
			Cooldown:                  m.Command.Cooldown,
			CooldownType:              m.Command.CooldownType,
			Enabled:                   m.Command.Enabled,
			Aliases:                   m.Command.Aliases,
			Description:               m.Command.Description,
			Visible:                   m.Command.Visible,
			ChannelID:                 m.Command.ChannelID,
			Default:                   m.Command.Default,
			DefaultName:               m.Command.DefaultName,
			Module:                    m.Command.Module,
			IsReply:                   m.Command.IsReply,
			KeepResponsesOrder:        m.Command.KeepResponsesOrder,
			DeniedUsersIDS:            m.Command.DeniedUsersIDS,
			AllowedUsersIDS:           m.Command.AllowedUsersIDS,
			RolesIDS:                  m.Command.RolesIDS,
			OnlineOnly:                m.Command.OnlineOnly,
			CooldownRolesIDs:          m.Command.CooldownRolesIDs,
			EnabledCategories:         m.Command.EnabledCategories,
			RequiredWatchTime:         m.Command.RequiredWatchTime,
			RequiredMessages:          m.Command.RequiredMessages,
			RequiredUsedChannelPoints: m.Command.RequiredUsedChannelPoints,
			GroupID:                   m.Command.GroupID,
			ExpiresAt:                 m.Command.ExpiresAt,
			OfflineOnly:               m.Command.OfflineOnly,
		},
	}

	if m.Command.ExpiresType != nil {
		expire := entity.CommandExpireType(*m.Command.ExpiresType)
		e.Command.ExpiresType = &expire
	}

	if m.Group != nil {
		e.Group = &entity.CommandGroup{
			ID:        m.Group.ID,
			ChannelID: m.Group.ChannelID,
			Name:      m.Group.Name,
			Color:     m.Group.Color,
		}
	}

	responses := make([]entity.CommandResponse, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses, entity.CommandResponse{
				ID:                r.ID,
				Text:              r.Text,
				CommandID:         r.CommandID,
				Order:             r.Order,
				TwitchCategoryIDs: r.TwitchCategoryIDs,
			},
		)
	}

	e.Responses = responses

	return e
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.CommandWithGroupAndResponses,
	error,
) {
	cmds, err := c.commandsWithGroupsAndResponsesRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.CommandWithGroupAndResponses, 0, len(cmds))
	for _, cmd := range cmds {
		entities = append(entities, c.mapToEntity(cmd))
	}

	slices.SortFunc(
		entities, func(a, b entity.CommandWithGroupAndResponses) int {
			return strings.Compare(a.Command.Name, b.Command.Name)
		},
	)

	return entities, nil
}
