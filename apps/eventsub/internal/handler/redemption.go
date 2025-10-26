package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistmodel "github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	channelredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"

	"github.com/google/uuid"
	model "github.com/twirapp/twir/libs/gomodels"
)

type userForIncrementUsedEmotes struct {
	userId    string
	channelId string
	cost      int
}

func (c *Handler) handleChannelPointsRewardRedemptionAddBatched(
	ctx context.Context,
	data []eventsub.ChannelPointsCustomRewardRedemptionAddEvent,
) {
	itemsForHistoryCreate := make([]channelredemptionshistory.CreateInput, len(data))
	itemsForEventsCreate := make([]channelseventslist.CreateInput, len(data))
	usersForIncrementUsedEmotes := make(map[string]*userForIncrementUsedEmotes)

	ctxWithoutCancel := context.WithoutCancel(ctx)

	for i, event := range data {
		c.logger.Info(
			"channel points reward redemption add",
			slog.String("reward", event.Reward.Title),
			slog.String("userName", event.UserLogin),
			slog.String("userId", event.UserId),
			slog.String("channelName", event.BroadcasterUserLogin),
			slog.String("channelId", event.BroadcasterUserId),
		)

		itemsForHistoryCreate[i] = channelredemptionshistory.CreateInput{
			ChannelID:    event.BroadcasterUserId,
			UserID:       event.UserId,
			RewardID:     uuid.MustParse(event.Reward.Id),
			RewardTitle:  event.Reward.Title,
			RewardPrompt: &event.UserInput,
			RewardCost:   event.Reward.Cost,
		}

		itemsForEventsCreate[i] = channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
			Type:      channelseventslistmodel.ChannelEventListItemTypeRedemptionCreated,
			Data: &channelseventslistmodel.ChannelsEventsListItemData{
				RedemptionInput:           event.UserInput,
				RedemptionTitle:           event.Reward.Title,
				RedemptionUserName:        event.UserLogin,
				RedemptionUserDisplayName: event.UserName,
				RedemptionCost:            strconv.Itoa(event.Reward.Cost),
			},
		}

		err := c.twirBus.Events.RedemptionCreated.Publish(
			ctx,
			events.RedemptionCreatedMessage{
				ID: event.Reward.Id,
				BaseInfo: events.BaseInfo{
					ChannelID:   event.BroadcasterUserId,
					ChannelName: event.BroadcasterUserLogin,
				},
				UserID:          event.UserId,
				UserName:        event.UserLogin,
				UserDisplayName: event.UserName,
				RewardName:      event.Reward.Title,
				RewardCost:      strconv.Itoa(event.Reward.Cost),
				Input:           lo.If(event.UserInput != "", &event.UserInput).Else(nil),
			},
		)
		if err != nil {
			c.logger.Error(err.Error(), slog.Any("err", err))
		}

		if _, ok := usersForIncrementUsedEmotes[event.BroadcasterUserId+event.UserId]; ok {
			usersForIncrementUsedEmotes[event.BroadcasterUserId+event.UserId].cost += event.Reward.Cost
		} else {
			usersForIncrementUsedEmotes[event.BroadcasterUserId+event.UserId] = &userForIncrementUsedEmotes{
				userId:    event.UserId,
				channelId: event.BroadcasterUserId,
				cost:      event.Reward.Cost,
			}
		}

		// youtube song requests

		go func() {
			e := c.handleYoutubeSongRequests(ctxWithoutCancel, &event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("e", err))
			}
		}()

		go func() {
			e := c.handleAlerts(ctxWithoutCancel, &event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("e", err))
			}
		}()

		go func() {
			e := c.handleRewardsSevenTvEmote(ctxWithoutCancel, &event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("err", e))
			}
		}()

		go func() {
			if redemptionAddErr := c.twirBus.RedemptionAdd.Publish(
				ctx,
				twitch.ActivatedRedemption{
					ID:                   event.Id,
					BroadcasterUserID:    event.BroadcasterUserId,
					BroadcasterUserLogin: event.BroadcasterUserLogin,
					BroadcasterUserName:  event.BroadcasterUserName,
					UserID:               event.UserId,
					UserLogin:            event.UserLogin,
					UserName:             event.UserName,
					UserInput:            event.UserInput,
					Status:               string(event.Status),
					RedeemedAt:           time.Now(),
					Reward: twitch.ActivatedRedemptionReward{
						ID:     event.Reward.Id,
						Title:  event.Reward.Title,
						Prompt: event.Reward.Prompt,
						Cost:   event.Reward.Cost,
					},
				},
			); redemptionAddErr != nil {
				c.logger.Error(redemptionAddErr.Error(), slog.Any("err", redemptionAddErr))
			}
		}()
	}

	go func() {
		for _, userForIncrement := range usersForIncrementUsedEmotes {
			err := c.countUserChannelPoints(
				ctxWithoutCancel,
				userForIncrement.userId,
				userForIncrement.channelId,
				userForIncrement.cost,
			)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}

		if len(itemsForHistoryCreate) > 0 {
			err := c.redemptionsHistoryRepository.CreateMany(ctxWithoutCancel, itemsForHistoryCreate)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}

		if len(itemsForEventsCreate) > 0 {
			err := c.eventsListRepository.CreateMany(ctxWithoutCancel, itemsForEventsCreate)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}
	}()
}

func (c *Handler) HandleChannelPointsRewardRedemptionAdd(
	ctx context.Context,
	event eventsub.ChannelPointsCustomRewardRedemptionAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.redemptionsBatcher.Add(event)
}

func (c *Handler) HandleChannelPointsRewardRedemptionUpdate(
	ctx context.Context,
	event eventsub.ChannelPointsCustomRewardRedemptionUpdateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	if event.Status != eventsub.CustomRewardRedemptionStatusCanceled {
		return
	}

	userStats := &model.UsersStats{}
	err := c.gorm.WithContext(ctx).Where(`"userId" = ?`, event.UserId).Find(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
	if userStats.ID == "" {
		return
	}
	userStats.UsedChannelPoints -= int64(event.Reward.Cost)
	err = c.gorm.WithContext(ctx).Save(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) countUserChannelPoints(
	ctx context.Context,
	userId, channelId string,
	count int,
) error {
	user := &model.Users{}
	err := c.gorm.
		WithContext(ctx).
		Where("id = ?", userId).
		Preload("Stats", `"channelId" = ?`, channelId).
		First(user).Error
	if err != nil {
		return err
	}

	if user.ID == "" {
		user = &model.Users{
			ID:         "",
			TokenID:    sql.NullString{},
			IsBotAdmin: false,
			ApiKey:     uuid.New().String(),
			Stats: &model.UsersStats{
				ID:                uuid.New().String(),
				UserID:            userId,
				ChannelID:         channelId,
				Messages:          0,
				Watched:           0,
				UsedChannelPoints: int64(count),
				Emotes:            0,
			},
		}

		err = c.gorm.Error
		if err != nil {
			return err
		}
	}

	if user.Stats != nil {
		user.Stats.UsedChannelPoints += int64(count)
		err = c.gorm.WithContext(ctx).Save(user.Stats).Error
		if err != nil {
			return err
		}
	} else {
		user.Stats = &model.UsersStats{
			ID:                uuid.New().String(),
			UserID:            userId,
			ChannelID:         channelId,
			Messages:          0,
			Watched:           0,
			UsedChannelPoints: int64(count),
			Emotes:            0,
		}
		err = c.gorm.WithContext(ctx).Create(user.Stats).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Handler) handleYoutubeSongRequests(
	ctx context.Context,
	event *eventsub.ChannelPointsCustomRewardRedemptionAddEvent,
) error {
	if event.UserInput == "" {
		return nil
	}

	entity, err := c.channelSongRequestsSettingsCache.Get(
		ctx,
		event.BroadcasterUserId,
	)
	if err != nil {
		return err
	}
	if entity.ID == "" || !entity.Enabled || event.Reward.Id != entity.ChannelPointsRewardID.String {
		return nil
	}

	var foundCommand *commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses
	commands, err := c.commandsCache.Get(ctx, event.BroadcasterUserId)
	if err != nil {
		return err
	}

	for _, command := range commands {
		if command.DefaultName != nil && *command.DefaultName == "sr" && command.Enabled {
			foundCommand = &command
			break
		}
	}

	if foundCommand == nil {
		return nil
	}

	res, err := c.twirBus.Parser.GetCommandResponse.Request(
		ctx,
		twitch.TwitchChatMessage{
			BroadcasterUserId:    event.BroadcasterUserId,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ChatterUserId:        event.UserId,
			ChatterUserName:      event.UserName,
			ChatterUserLogin:     event.UserLogin,
			MessageId:            event.Id,
			Message: &twitch.ChatMessageMessage{
				Text: fmt.Sprintf("!%s %s", foundCommand.Name, event.UserInput),
			},
		},
	)
	if err != nil {
		return err
	}

	if len(res.Data.Responses) == 0 {
		return nil
	}

	for _, response := range res.Data.Responses {
		c.twirBus.Bots.SendMessage.Publish(
			ctx,
			bots.SendMessageRequest{
				ChannelId:      event.BroadcasterUserId,
				ChannelName:    &event.BroadcasterUserLogin,
				Message:        fmt.Sprintf("@%s %s", event.UserLogin, response),
				SkipRateLimits: true,
			},
		)
	}

	return nil
}

func (c *Handler) handleAlerts(
	ctx context.Context,
	event *eventsub.ChannelPointsCustomRewardRedemptionAddEvent,
) error {
	alerts, err := c.alertsCache.Get(ctx, event.BroadcasterUserId)
	if err != nil {
		return err
	}

	var foundAlertId uuid.UUID
	for _, alert := range alerts {
		if slices.Contains(alert.RewardIDS, event.Reward.Id) {
			foundAlertId = alert.ID
			break
		}
	}
	if foundAlertId == uuid.Nil {
		return nil
	}

	_, err = c.websocketsGrpc.TriggerAlert(
		ctx,
		&websockets.TriggerAlertRequest{
			ChannelId: event.BroadcasterUserId,
			AlertId:   foundAlertId.String(),
		},
	)

	return err
}
