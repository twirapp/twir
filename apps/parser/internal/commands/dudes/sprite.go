package dudes

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/websockets"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/types/types/overlays"
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
	SkipToxicityCheck: true,
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
				result.Result = []string{i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dudes.Info.Sprite.
						SetVars(locales.KeysCommandsDudesInfoSpriteVars{DudeSprite: *entity.DudeColor}),
				)}
				return &result, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dudes.Info.SpriteRequired.
						SetVars(locales.KeysCommandsDudesInfoSpriteRequiredVars{AvailableSprites: strings.Join(availableSpritesStr, ", ")}),
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
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dudes.Errors.SpriteInvalid.
						SetVars(locales.KeysCommandsDudesErrorsSpriteInvalidVars{AvailableSprites: strings.Join(availableSpritesStr, ", ")}),
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
			ctx,
			websockets.DudesChangeUserSettingsRequest{
				ChannelID: parseCtx.Channel.ID,
				UserID:    parseCtx.Sender.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dudes.Errors.SpriteCannotTrigger),
				Err:     err,
			}
		}

		result.Result = []string{i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Dudes.Info.SpriteChanged.
				SetVars(locales.KeysCommandsDudesInfoSpriteChangedVars{DudeSprite: sprite.String()}),
		)}
		return &result, nil
	},
}
