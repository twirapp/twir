package dudes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/mazznoer/csscolorparser"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

var Color = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes color",
		Description: null.StringFrom("Change the color of user in the dudes overlay"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
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

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			if entity.UserID != "" && entity.DudeColor != nil {
				result.Result = []string{fmt.Sprintf("Your color is %s", *entity.DudeColor)}
				return &result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: "color is required",
			}
		}

		color, err := csscolorparser.Parse(*parseCtx.Text)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "invalid color",
				Err:     err,
			}
		}

		if entity.UserID == "" {
			entity.ID = uuid.New()
			entity.ChannelID = parseCtx.Channel.ID
			entity.UserID = parseCtx.Sender.ID
		}

		entity.DudeColor = lo.ToPtr(color.HexString())
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&entity).Error; err != nil {
			return nil, err
		}

		err = parseCtx.Services.Bus.WebsocketsDudesUserSettings.Publish(
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

		result.Result = []string{fmt.Sprintf("Color changed to %s", color.HexString())}
		return &result, nil
	},
}
