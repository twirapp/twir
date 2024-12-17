package commands_with_groups_and_responses

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CommandsRepository commands_with_groups_and_responses.Repository
}

func New(opts Opts) *Service {
	return &Service{
		commandsRepository: opts.CommandsRepository,
	}
}

type Service struct {
	commandsRepository commands_with_groups_and_responses.Repository
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
	commands, err := c.commandsRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.CommandWithGroupAndResponses, 0, len(commands))
	for _, cmd := range commands {
		entities = append(entities, c.mapToEntity(cmd))
	}

	return entities, nil
}
