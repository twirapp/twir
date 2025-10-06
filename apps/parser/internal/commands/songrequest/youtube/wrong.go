package sr_youtube

import (
	"context"
	"time"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
)

const (
	songSkipArgName = "number"
)

var WrongCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr wrong",
		Description: null.StringFrom("Delete wrong song from queue"),
		Module:      "SONGS",
		IsReply:     true,
		Visible:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     songSkipArgName,
			Optional: true,
			Min:      lo.ToPtr(1),
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		var songs []*model.RequestedSong
		err := parseCtx.Services.Gorm.WithContext(ctx).
			Where(
				`"channelId" = ? AND "orderedById" = ? AND "deletedAt" IS NULL`,
				parseCtx.Channel.ID,
				parseCtx.Sender.ID,
			).
			Limit(5).
			Order(`"createdAt" desc`).
			Find(&songs).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.GetSongsFromQueue),
				Err:     err,
			}
		}

		if len(songs) == 0 {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Info.NoRequestedSongs))
			return result, nil
		}

		number := 1
		songSkipArg := parseCtx.ArgsParser.Get(songSkipArgName)
		if songSkipArg != nil {
			number = songSkipArg.Int()
		}

		if number > len(songs)+1 {
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Songrequest.Info.OnlyCountSongs.
						SetVars(locales.KeysCommandsSongrequestInfoOnlyCountSongsVars{SongsCount: len(songs)}),
				),
			)
			return result, nil
		}

		choosedSong := songs[number-1]
		choosedSong.DeletedAt = lo.ToPtr(time.Now().UTC())
		err = parseCtx.Services.Gorm.WithContext(ctx).Updates(&choosedSong).Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.UpdateSong),
				Err:     err,
			}
		}

		_, err = parseCtx.Services.GrpcClients.WebSockets.YoutubeRemoveSongToQueue(
			ctx,
			&websockets.YoutubeRemoveSongFromQueueRequest{
				ChannelId: parseCtx.Channel.ID,
				EntityId:  choosedSong.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Songrequest.Errors.RemoveSongFromQueue),
				Err:     err,
			}
		}

		result.Result = append(
			result.Result,
			i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Songrequest.Info.Delete.
					SetVars(locales.KeysCommandsSongrequestInfoDeleteVars{SongTitle: choosedSong.Title}),
			),
		)

		return result, nil
	},
}
