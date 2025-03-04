package http_public

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
)

type publicCommandsOutput struct {
	Body []commandDto `json:"body"`
}

type commandDto struct {
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Module      string               `json:"module"`
	Group       *string              `json:"group"`
	Responses   []commandDtoResponse `json:"responses"`
}

type commandDtoResponse struct {
	Text string `json:"text"`
}

func (p *Public) HandleChannelCommandsGet(
	ctx context.Context,
	channelId string,
) (*publicCommandsOutput, error) {
	channel, err := p.channelsService.GetByID(ctx, channelId)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return nil, huma.Error404NotFound("channel not found")
		}

		return nil, huma.Error500InternalServerError("internal server error")
	}

	if !channel.IsEnabled {
		return nil, huma.Error404NotFound("channel not found")
	}

	// TODO: it must use generic cacher with repository instead of gorm cacher
	commands, err := p.cachedCommands.Get(ctx, channelId)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}

	commands = lo.Filter(
		commands, func(item model.ChannelsCommands, index int) bool {
			return item.Enabled && item.Visible
		},
	)

	mappedCommands := make([]commandDto, 0, len(commands))

	for _, command := range commands {
		mappedCmd := commandDto{
			Name:        command.Name,
			Description: command.Description.Ptr(),
			Module:      command.Module,
			Group:       nil,
			Responses:   make([]commandDtoResponse, 0, len(command.Responses)),
		}

		if command.Group != nil {
			mappedCmd.Group = &command.Group.Name
		}

		for _, response := range command.Responses {
			mappedCmd.Responses = append(
				mappedCmd.Responses,
				commandDtoResponse{
					Text: response.Text.String,
				},
			)
		}

		mappedCommands = append(mappedCommands, mappedCmd)
	}

	return &publicCommandsOutput{Body: mappedCommands}, nil
}
