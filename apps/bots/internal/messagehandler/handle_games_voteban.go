package messagehandler

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	channelsgamesvotebanprogressstate "github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state"
	progressstatemodel "github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state/model"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleGamesVoteban(ctx context.Context, msg handleMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	mu := c.votebanLock.NewMutex("bots:voteban_handle_message:" + msg.BroadcasterUserId)
	mu.Lock()
	defer mu.Unlock()

	// Check if voting is not in progress
	voteExists, err := c.votebanProgressStateRepository.Exists(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	if !voteExists {
		return nil
	}

	// Check if user has already voted
	userVoted, err := c.votebanProgressStateRepository.UserHasVoted(
		ctx,
		msg.BroadcasterUserId,
		msg.ChatterUserId,
	)
	if err != nil {
		return err
	}

	if userVoted {
		return nil
	}

	// Get current vote state
	voteState, err := c.votebanProgressStateRepository.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		if errors.Is(err, channelsgamesvotebanprogressstate.ErrNotFound) {
			return nil
		}
		return err
	}

	// Get game settings
	gameEntity, err := c.channelsGamesVotebanCacher.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		if errors.Is(err, channelsgamesvoteban.ErrNotFound) {
			return nil
		}
		return err
	}

	if gameEntity.IsNil() {
		return nil
	}

	if !gameEntity.Enabled {
		return nil
	}

	splittedChatMessage := strings.Fields(msg.Message.Text)

	for _, word := range splittedChatMessage {
		if slices.Contains(gameEntity.ChatVotesWordsPositive, word) {
			voteState.TotalVotes++
			voteState.PositiveVotes++
			break
		} else if slices.Contains(gameEntity.ChatVotesWordsNegative, word) {
			voteState.TotalVotes++
			voteState.NegativeVotes++
			break
		}
	}

	if voteState.TotalVotes >= gameEntity.NeededVotes {
		isPositive := voteState.PositiveVotes > voteState.NegativeVotes

		var message string
		if isPositive {
			if voteState.TargetIsMod {
				message = gameEntity.BanMessageModerators
			} else {
				message = gameEntity.BanMessage
			}
		} else {
			if voteState.TargetIsMod {
				message = gameEntity.SurviveMessageModerators
			} else {
				message = gameEntity.SurviveMessage
			}
		}

		message = strings.ReplaceAll(message, "{targetUser}", voteState.TargetUserName)

		if isPositive {
			if err := c.twitchActions.Ban(
				ctx,
				twitchactions.BanOpts{
					Duration:       gameEntity.TimeoutSeconds,
					Reason:         message,
					BroadcasterID:  msg.BroadcasterUserId,
					UserID:         voteState.TargetUserID,
					ModeratorID:    msg.EnrichedData.DbChannel.BotID,
					IsModerator:    voteState.TargetIsMod,
					AddModAfterBan: true,
				},
			); err != nil {
				return err
			}
		}

		if err := c.twitchActions.SendMessage(
			ctx,
			twitchactions.SendMessageOpts{
				BroadcasterID: msg.BroadcasterUserId,
				SenderID:      msg.EnrichedData.DbChannel.BotID,
				Message:       message,
			},
		); err != nil {
			return err
		}

		// Delete vote state and clear user votes
		if err := c.votebanProgressStateRepository.Delete(ctx, msg.BroadcasterUserId); err != nil {
			return err
		}

		if err := c.votebanProgressStateRepository.ClearUserVotes(
			ctx,
			msg.BroadcasterUserId,
		); err != nil {
			return err
		}
	} else {
		// Update vote state
		if err := c.votebanProgressStateRepository.Update(
			ctx,
			msg.BroadcasterUserId,
			voteState,
		); err != nil {
			return err
		}
	}

	return nil
}

// mapVoteStateToModel converts internal VoteState to progressstatemodel.VoteState
func mapVoteStateToProgressModel(state progressstatemodel.VoteState) progressstatemodel.VoteState {
	return state
}
