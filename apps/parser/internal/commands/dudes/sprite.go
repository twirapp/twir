package dudes

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/types/types/overlays"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

const (
	spriteArgName = "sprite"
)

var Sprite = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes sprite",
		Description: null.StringFrom("Change sprite of dude"),
		RolesIDS:    pq.StringArray{},
		Module:      "DUDES",
		Visible:     true,
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: spriteArgName,
			OneOf: lo.Map(
				overlays.AllDudesSpriteEnumValues,
				func(item overlays.DudesSprite, _ int) string {
					return item.String()
				},
			),
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		entity := model.ChannelsOverlaysDudesUserSettings{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`channel_id = ? AND user_id = ?`, parseCtx.Channel.ID, parseCtx.Sender.ID).
			Find(&entity).Error; err != nil {
			return nil, err
		}

		result := types.CommandsHandlerResult{}

		availableSprites := overlays.AllDudesSpriteEnumValues
		availableSpritesStr := make([]string, len(availableSprites))
		for i, v := range availableSprites {
			availableSpritesStr[i] = v.String()
		}

		spriteArg := parseCtx.ArgsParser.Get(spriteArgName)

		if spriteArg == nil {
			if entity.UserID != "" && entity.DudeSprite != nil {
				result.Result = []string{fmt.Sprintf("Your sprite it %s", *entity.DudeSprite)}
				return &result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"sprite is required, available: %v",
					strings.Join(availableSpritesStr, ", "),
				),
			}
		}

		if entity.UserID == "" {
			entity.ID = uuid.New()
			entity.ChannelID = parseCtx.Channel.ID
			entity.UserID = parseCtx.Sender.ID
		}

		sprite := overlays.DudesSprite(spriteArg.String())
		if !sprite.IsValid() {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"invalid sprite, available: %v",
					strings.Join(availableSpritesStr, ", "),
				),
			}
		}

		entity.DudeSprite = lo.ToPtr(sprite.String())
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
				Message: "cannot trigger dudes sprite",
				Err:     err,
			}
		}

		result.Result = []string{fmt.Sprintf("Sprite changed to %s", sprite.String())}
		return &result, nil
	},
}
