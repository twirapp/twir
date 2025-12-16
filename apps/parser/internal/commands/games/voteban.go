package games

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/bus-core/bots"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	"github.com/twirapp/twir/libs/repositories/userswithstats"
)

const (
	votebanArgName = "user"
)

var Voteban = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "voteban",
		Description: null.StringFrom("Initiate voteban"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: votebanArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		mu := parseCtx.Services.RedSync.NewMutex(
			"parser:voteban:"+parseCtx.Channel.ID,
			redsync.WithExpiry(5*time.Second),
		)
		if err := mu.Lock(); err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotLock),
				Err:     err,
			}
		}
		defer mu.Unlock()

		entity, err := parseCtx.Services.ChannelsGamesVotebanRepo.GetByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			if errors.Is(err, channelsgamesvoteban.ErrNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Games.Errors.VotebanCannotFindSettings,
				),
				Err: err,
			}
		}

		if !entity.Enabled {
			return nil, nil
		}

		if len(parseCtx.Mentions) == 0 {
			return nil, nil
		}

		targetUser := parseCtx.Mentions[0]

		// Fetch channel to check BotID
		dbChannel := model.Channels{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"id" = ?`, parseCtx.Channel.ID).
			First(&dbChannel).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Games.Errors.VotebanCannotFindSettings,
				),
				Err: err,
			}
		}

		if targetUser.UserId == parseCtx.Channel.ID ||
			targetUser.UserId == dbChannel.BotID {
			return nil, nil
		}

		targerUserDbStats, err := parseCtx.Services.UsersWithStatsRepository.GetByUserAndChannelID(
			ctx, userswithstats.GetByUserAndChannelIDInput{
				UserID:    targetUser.UserId,
				ChannelID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		if targerUserDbStats.Stats == nil {
			return &types.CommandsHandlerResult{
				Result: []string{"user not found"},
			}, nil
		}

		if targerUserDbStats.Stats.IsMod && !entity.TimeoutModerators {
			return &types.CommandsHandlerResult{
				Result: []string{i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Games.Errors.VotebanCannotTimeoutModerator,
				)},
			}, nil
		}

		res, err := parseCtx.Services.Bus.Bots.VotebanRegister.Request(
			ctx, bots.VotebanRegisterRequest{
				Data:                 entity,
				TargerUser:           targetUser,
				InitiatorUserID:      parseCtx.Sender.ID,
				InitiatorUserLogin:   parseCtx.Sender.Name,
				InitiatorIsModerator: targerUserDbStats.Stats.IsMod,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Games.Errors.VotebanCannotCheckProgress,
				),
				Err: err,
			}
		}

		if res.Data.AlreadyInProgress {
			return &types.CommandsHandlerResult{
				Result: []string{i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Games.Info.VotebanInProgress,
				)},
			}, nil
		}

		initMessage := strings.ReplaceAll(
			entity.InitMessage,
			"{targetUser}",
			targetUser.UserName,
		)
		initMessage = strings.ReplaceAll(
			initMessage,
			"{positiveTexts}",
			strings.Join(entity.ChatVotesWordsPositive, " · "),
		)
		initMessage = strings.ReplaceAll(
			initMessage,
			"{negativeTexts}",
			strings.Join(entity.ChatVotesWordsNegative, " · "),
		)

		return &types.CommandsHandlerResult{
			Result: []string{initMessage},
		}, nil
	},
}
