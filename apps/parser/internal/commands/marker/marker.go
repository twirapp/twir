package marker

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

var Marker = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "marker",
		Description: null.StringFrom("Create a stream marker"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     "markerDescription",
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		dbChannel := &model.Channels{}
		if err := parseCtx.Services.Gorm.First(dbChannel, parseCtx.Channel.ID).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.CannotFindChannelDb,
				),
				Err: err,
			}
		}

		twitchClient, err := twitch.NewBotClientWithContext(
			ctx,
			dbChannel.BotID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Errors.Generic.BroadcasterClient,
				),
				Err: err,
			}
		}

		params := helix.CreateStreamMarkerParams{
			UserID:      dbChannel.ID,
			Description: "",
		}

		description := parseCtx.ArgsParser.Get("markerDescription")
		if description != nil {
			params.Description = description.String()
		}

		resp, err := twitchClient.CreateStreamMarker(&params)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Marker.Errors.CannotCreateMarker.SetVars(locales.KeysCommandsMarkerErrorsCannotCreateMarkerVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}
		if resp.StatusCode == 403 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Marker.Errors.CannotCreateMarker.SetVars(locales.KeysCommandsMarkerErrorsCannotCreateMarkerVars{Reason: "insufficient permissions"}),
				),
				Err: errors.New("insufficient permissions"),
			}
		}

		if resp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Marker.Errors.CannotCreateMarker.SetVars(locales.KeysCommandsMarkerErrorsCannotCreateMarkerVars{Reason: resp.ErrorMessage}),
				),
				Err: err,
			}
		}

		if len(resp.Data.CreateStreamMarkers) == 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Marker.Errors.CannotCreateMarker,
				),
				Err: errors.New("empty stream marker response"),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Marker.Success.MarkerCreated,
				),
			},
		}, nil
	},
}
