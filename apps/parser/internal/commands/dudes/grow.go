package dudes

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

var Grow = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes grow",
		Description: null.StringFrom("Increase the size of user in the dudes overlay"),
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

		var color string
		if entity.UserID != "" {
			color = *entity.DudeColor
		} else {
			color = parseCtx.Sender.Color
		}

		err := parseCtx.Services.Bus.Websocket.DudesGrow.Publish(
			websockets.DudesGrowRequest{
				ChannelID:       parseCtx.Channel.ID,
				UserID:          parseCtx.Sender.ID,
				UserDisplayName: parseCtx.Sender.DisplayName,
				UserName:        parseCtx.Sender.Name,
				Color:           color,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger dudes grow",
				Err:     err,
			}
		}

		return &result, nil
	},
}
