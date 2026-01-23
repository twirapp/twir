package commands_with_groups_and_responses

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	commandsservice "github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/libs/audit"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	"github.com/twirapp/twir/libs/repositories/command_role_cooldown"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TrmManager                               trm.Manager
	CommandsRepository                       commands.Repository
	CommandsWithGroupsAndResponsesRepository commands_with_groups_and_responses.Repository
	ResponsesRepository                      commands_response.Repository
	CommandRoleCooldownRepository            command_role_cooldown.Repository
	CommandsService                          *commandsservice.Service
	Logger                                   *slog.Logger
	AuditRecorder                            audit.Recorder
	CachedCommandsClient                     *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
}

func New(opts Opts) *Service {
	return &Service{
		trmManager:                               opts.TrmManager,
		commandsWithGroupsAndResponsesRepository: opts.CommandsWithGroupsAndResponsesRepository,
		responsesRepository:                      opts.ResponsesRepository,
		commandsRepository:                       opts.CommandsRepository,
		commandRoleCooldownRepository:            opts.CommandRoleCooldownRepository,
		logger:                                   opts.Logger,
		auditRecorder:                            opts.AuditRecorder,
		commandsService:                          opts.CommandsService,
		cachedCommandsClient:                     opts.CachedCommandsClient,
	}
}

type Service struct {
	trmManager                               trm.Manager
	commandsRepository                       commands.Repository
	commandsWithGroupsAndResponsesRepository commands_with_groups_and_responses.Repository
	responsesRepository                      commands_response.Repository
	commandRoleCooldownRepository            command_role_cooldown.Repository

	commandsService      *commandsservice.Service
	logger               *slog.Logger
	auditRecorder        audit.Recorder
	cachedCommandsClient *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
}

func (c *Service) mapToEntity(m model.CommandWithGroupAndResponses) commandwithrelationentity.CommandWithGroupAndResponses {
	e := commandwithrelationentity.CommandWithGroupAndResponses{
		Command: commandwithrelationentity.Command{
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
		expire := commandwithrelationentity.CommandExpireType(*m.Command.ExpiresType)
		e.Command.ExpiresType = &expire
	}

	if m.Group != nil {
		e.Group = &commandwithrelationentity.CommandGroup{
			ID:        m.Group.ID,
			ChannelID: m.Group.ChannelID,
			Name:      m.Group.Name,
			Color:     m.Group.Color,
		}
	}

	responses := make([]commandwithrelationentity.CommandResponse, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses, commandwithrelationentity.CommandResponse{
				ID:                r.ID,
				Text:              r.Text,
				CommandID:         r.CommandID,
				Order:             r.Order,
				TwitchCategoryIDs: r.TwitchCategoryIDs,
				OnlineOnly:        r.OnlineOnly,
				OfflineOnly:       r.OfflineOnly,
			},
		)
	}

	e.Responses = responses

	return e
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]commandwithrelationentity.CommandWithGroupAndResponses,
	error,
) {
	cmds, err := c.commandsWithGroupsAndResponsesRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]commandwithrelationentity.CommandWithGroupAndResponses, 0, len(cmds))
	for _, cmd := range cmds {
		entities = append(entities, c.mapToEntity(cmd))
	}

	slices.SortFunc(
		entities, func(a, b commandwithrelationentity.CommandWithGroupAndResponses) int {
			return strings.Compare(a.Command.Name, b.Command.Name)
		},
	)

	return entities, nil
}
