package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/libs/entities/platform"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

type publicCommandsOutput struct {
	Body []commandV2Dto `json:"body"`
}

type commandV2Dto struct {
	Name        string                   `json:"name"`
	Description *string                  `json:"description"`
	Module      string                   `json:"module"`
	Group       *string                  `json:"group"`
	Responses   []kickCommandDtoResponse `json:"responses"`
}

type kickCommandDtoResponse struct {
	Text string `json:"text"`
}

func (p *Public) getChannelByPlatformAndID(ctx context.Context, inputPlatform platform.Platform, channelID string) (entity.Channel, error) {
	var (
		channel      entity.Channel
		channelError error
	)

	if inputPlatform == platform.PlatformKick {
		channel, channelError = p.channelsService.GetByKickPlatformID(ctx, channelID)
	} else if inputPlatform == platform.PlatformTwitch {
		channel, channelError = p.channelsService.GetByTwitchPlatformID(ctx, channelID)
	} else {
		return entity.ChannelNil, fmt.Errorf("invalid platform")
	}

	if channelError != nil {
		if errors.Is(channelError, channels.ErrNotFound) {
			return entity.ChannelNil, huma.Error404NotFound("channel not found")
		}

		return entity.ChannelNil, huma.Error500InternalServerError("internal server error")
	}

	return channel, nil
}

func (p *Public) getChannelCommands(ctx context.Context, channelID uuid.UUID) (*publicCommandsOutput, error) {
	commands, err := p.cachedCommands.Get(ctx, channelID.String())
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}

	commands = lo.Filter(
		commands,
		func(item commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses, index int) bool {
			return item.Enabled && item.Visible
		},
	)

	mappedCommands := make([]commandV2Dto, 0, len(commands))

	for _, command := range commands {
		mappedCmd := commandV2Dto{
			Name:        command.Name,
			Description: command.Description,
			Module:      command.Module,
			Group:       nil,
			Responses:   make([]kickCommandDtoResponse, 0, len(command.Responses)),
		}

		if command.Group != nil {
			mappedCmd.Group = &command.Group.Name
		}

		for _, response := range command.Responses {
			var text string
			if response.Text != nil {
				text = *response.Text
			}

			mappedCmd.Responses = append(
				mappedCmd.Responses,
				kickCommandDtoResponse{
					Text: text,
				},
			)
		}

		mappedCommands = append(mappedCommands, mappedCmd)
	}

	return &publicCommandsOutput{Body: mappedCommands}, nil
}

func (p *Public) handleGetChannelByPlatformCommandsGet(ctx context.Context, inputPlatform platform.Platform, id string) (*publicCommandsOutput, error) {
	channel, err := p.getChannelByPlatformAndID(ctx, inputPlatform, id)
	if err != nil {
		return nil, err
	}

	commands, err := p.getChannelCommands(ctx, channel.ID)
	if err != nil {
		return nil, err
	}

	return commands, nil
}
