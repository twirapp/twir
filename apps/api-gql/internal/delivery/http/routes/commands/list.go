package commands

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	"go.uber.org/fx"
)

type listRouteRequestDto struct {
	ChannelID string `path:"channel_id" required:"true" description:"The ID of the channel to get commands for"`
	Module    string `query:"module" description:"Filter commands by module name" example:"CUSTOM"`
}

var _ httpbase.Route[*listRouteRequestDto, *httpbase.BaseOutputJson[[]commandResponseDto]] = (*listById)(nil)

type ListByIdOpts struct {
	fx.In

	Service        *commands_with_groups_and_responses.Service
	ChannelService *channels.Service
}

type listById struct {
	service        *commands_with_groups_and_responses.Service
	channelService *channels.Service
}

func newListById(opts ListByIdOpts) *listById {
	return &listById{
		service:        opts.Service,
		channelService: opts.ChannelService,
	}
}

func (l *listById) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "commands-get-list-by-channel-id",
		Method:      http.MethodGet,
		Path:        "/v1/channels/{channel_id}/commands",
		Tags:        []string{"Commands"},
		Summary:     "Get commands list by channel id",
	}
}

func (l *listById) Handler(
	ctx context.Context,
	input *listRouteRequestDto,
) (*httpbase.BaseOutputJson[[]commandResponseDto], error) {
	_, err := l.channelService.GetByID(ctx, input.ChannelID)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Channel not found")
		}

		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get channel", err)
	}

	commands, err := l.service.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get commands", err)
	}

	output := make([]commandResponseDto, 0, len(commands))
	for _, cmd := range commands {
		if input.Module != "" && cmd.Command.Module != input.Module {
			continue
		}

		var group *commandGroupResponseDto
		var expire *Expire

		if cmd.Group != nil && *cmd.Group != commandwithrelationentity.CommandGroupNil {
			group = &commandGroupResponseDto{
				ID:    cmd.Group.ID,
				Name:  cmd.Group.Name,
				Color: cmd.Group.Color,
			}
		}

		if cmd.Command.ExpiresAt != nil {
			expire = &Expire{
				ExpiresAt:   *cmd.Command.ExpiresAt,
				ExpiresType: string(*cmd.Command.ExpiresType),
			}
		}

		responses := make([]commandResponsesResponseDto, 0, len(cmd.Responses))
		for _, resp := range cmd.Responses {
			var text string
			if resp.Text != nil {
				text = *resp.Text
			}

			responses = append(
				responses,
				commandResponsesResponseDto{
					ID:                resp.ID,
					Text:              text,
					Order:             resp.Order,
					TwitchCategoryIDs: resp.TwitchCategoryIDs,
					OnlineOnly:        resp.OnlineOnly,
					OfflineOnly:       resp.OfflineOnly,
				},
			)
		}

		rolesCooldowns := make([]commandRoleCooldownResponseDto, 0, len(cmd.RolesCooldowns))
		for _, c := range cmd.RolesCooldowns {
			rolesCooldowns = append(rolesCooldowns, commandRoleCooldownResponseDto{
				RoleID:   c.RoleID,
				Cooldown: c.Cooldown,
			})
		}

		dto := commandResponseDto{
			ID:                        cmd.Command.ID,
			Name:                      cmd.Command.Name,
			Cooldown:                  cmd.Command.Cooldown,
			CooldownType:              cmd.Command.CooldownType,
			Enabled:                   cmd.Command.Enabled,
			Aliases:                   append([]string{}, cmd.Command.Aliases...),
			Description:               cmd.Command.Description,
			Visible:                   cmd.Command.Visible,
			IsDefault:                 cmd.Command.Default,
			DefaultName:               cmd.Command.DefaultName,
			Module:                    cmd.Command.Module,
			IsReply:                   cmd.Command.IsReply,
			KeepResponsesOrder:        cmd.Command.KeepResponsesOrder,
			DeniedUsersIDS:            append([]string{}, cmd.Command.DeniedUsersIDS...),
			AllowedUsersIDS:           append([]string{}, cmd.Command.AllowedUsersIDS...),
			RolesIDS:                  append([]uuid.UUID{}, cmd.Command.RolesIDS...),
			OnlineOnly:                cmd.Command.OnlineOnly,
			OfflineOnly:               cmd.Command.OfflineOnly,
			EnabledCategories:         append([]string{}, cmd.Command.EnabledCategories...),
			RequiredWatchTime:         cmd.Command.RequiredWatchTime,
			RequiredMessages:          cmd.Command.RequiredMessages,
			RequiredUsedChannelPoints: cmd.Command.RequiredUsedChannelPoints,
			Expire:                    expire,
			Responses:                 responses,
			Group:                     group,
			RolesCooldowns:            rolesCooldowns,
		}

		output = append(output, dto)
	}

	return httpbase.CreateBaseOutputJson(output), nil
}

func (l *listById) Register(api huma.API) {
	huma.Register(api, l.GetMeta(), l.Handler)
}
