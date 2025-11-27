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
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	votebanmodel "github.com/twirapp/twir/libs/repositories/channels_games_voteban/model"
	progressstatemodel "github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state/model"
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

		entity, err := parseCtx.Services.ChannelsGamesVotebanRepo.GetByChannelID(ctx, parseCtx.Channel.ID)
		if err != nil {
			if errors.Is(err, channelsgamesvoteban.ErrNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotFindSettings),
				Err:     err,
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
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotFindSettings),
				Err:     err,
			}
		}

		if targetUser.UserId == parseCtx.Channel.ID ||
			targetUser.UserId == dbChannel.BotID {
			return nil, nil
		}

		if entity.VotingMode == votebanmodel.VotingModeChat {
			// Check if vote is already in progress
			voteInProgress, err := parseCtx.Services.ChannelsGamesVotebanProgressState.Exists(ctx, parseCtx.Channel.ID)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotCheckProgress),
					Err:     err,
				}
			}

			if voteInProgress {
				return &types.CommandsHandlerResult{
					Result: []string{i18n.GetCtx(ctx, locales.Translations.Commands.Games.Info.VotebanInProgress)},
				}, nil
			}

			targetUserStatsEntity := model.UsersStats{}
			if err := parseCtx.Services.Gorm.
				WithContext(ctx).
				Where(`"userId" = ? AND "channelId" = ?`, targetUser.UserId, parseCtx.Channel.ID).
				First(&targetUserStatsEntity).Error; err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotFindUser),
					Err:     err,
				}
			}

			voteDuration := time.Second * time.Duration(entity.VoteDuration)

			// Create vote state using repository
			if err := parseCtx.Services.ChannelsGamesVotebanProgressState.Create(
				ctx,
				parseCtx.Channel.ID,
				progressstatemodel.VoteState{
					TargetUserID:   targetUser.UserId,
					TargetUserName: targetUser.UserName,
					TargetIsMod:    targetUserStatsEntity.IsMod,
					TotalVotes:     1,
					PositiveVotes:  1,
					NegativeVotes:  0,
				},
				voteDuration,
			); err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Games.Errors.VotebanCannotSetVote),
					Err:     err,
				}
			}

			// Mark user as voted
			if err := parseCtx.Services.ChannelsGamesVotebanProgressState.MarkUserVoted(
				ctx,
				parseCtx.Channel.ID,
				parseCtx.Sender.ID,
				voteDuration,
			); err != nil {
				// Cleanup on error
				parseCtx.Services.ChannelsGamesVotebanProgressState.Delete(ctx, parseCtx.Channel.ID)
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
		}

		return &types.CommandsHandlerResult{
			Result: []string{},
		}, nil
	},
}
