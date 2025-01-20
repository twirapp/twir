package dudes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/mazznoer/csscolorparser"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

const (
	colorArgName = "color"
)

var Color = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes color",
		Description: null.StringFrom("Change the color of user in the dudes overlay"),
		Module:      "DUDES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: colorArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		entity := model.ChannelsOverlaysDudesUserSettings{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`channel_id = ? AND user_id = ?`, parseCtx.Channel.ID, parseCtx.Sender.ID).
			Find(&entity).Error; err != nil {
			return nil, err
		}

		text := parseCtx.ArgsParser.Get(colorArgName).String()

		if text == "" {
			if entity.UserID != "" && entity.DudeColor != nil {
				result.Result = []string{fmt.Sprintf("Your color is %s", *entity.DudeColor)}
				return &result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: "color is required",
			}
		}

		var color *string
		if text == "reset" {
			color = nil
		} else {
			parsedColor, err := csscolorparser.Parse(text)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "invalid color",
					Err:     err,
				}
			}

			color = lo.ToPtr(parsedColor.HexString())
		}

		if entity.UserID == "" {
			entity.ID = uuid.New()
			entity.ChannelID = parseCtx.Channel.ID
			entity.UserID = parseCtx.Sender.ID
		}

		entity.DudeColor = color
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&entity).Error; err != nil {
			return nil, err
		}

		err := parseCtx.Services.Bus.Websocket.DudesUserSettings.Publish(
			websockets.DudesChangeUserSettingsRequest{
				ChannelID: parseCtx.Channel.ID,
				UserID:    parseCtx.Sender.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger dudes color",
				Err:     err,
			}
		}

		if color == nil {
			result.Result = []string{"Color reset to default"}
			return &result, nil
		}

		result.Result = []string{fmt.Sprintf("Color changed to %s", *color)}
		return &result, nil
	},
}
