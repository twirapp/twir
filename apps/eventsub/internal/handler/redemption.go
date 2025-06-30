package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistmodel "github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	channelredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

type userForIncrementUsedEmotes struct {
	userId    string
	channelId string
	cost      int
}

func (c *Handler) handleChannelPointsRewardRedemptionAddBatched(
	ctx context.Context,
	data []eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) {
	itemsForHistoryCreate := make([]channelredemptionshistory.CreateInput, len(data))
	itemsForEventsCreate := make([]channelseventslist.CreateInput, len(data))
	usersForIncrementUsedEmotes := make(map[string]*userForIncrementUsedEmotes)

	for i, event := range data {
		c.logger.Info(
			"channel points reward redemption add",
			slog.String("reward", event.Reward.Title),
			slog.String("userName", event.UserLogin),
			slog.String("userId", event.UserID),
			slog.String("channelName", event.BroadcasterUserLogin),
			slog.String("channelId", event.BroadcasterUserID),
		)

		itemsForHistoryCreate[i] = channelredemptionshistory.CreateInput{
			ChannelID:    event.BroadcasterUserID,
			UserID:       event.UserID,
			RewardID:     uuid.MustParse(event.Reward.ID),
			RewardTitle:  event.Reward.Title,
			RewardPrompt: &event.UserInput,
			RewardCost:   event.Reward.Cost,
		}

		itemsForEventsCreate[i] = channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
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
			events.RedemptionCreatedMessage{
				ID: event.Reward.ID,
				BaseInfo: events.BaseInfo{
					ChannelID:   event.BroadcasterUserID,
					ChannelName: event.BroadcasterUserLogin,
				},
				UserID:          event.UserID,
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

		if _, ok := usersForIncrementUsedEmotes[event.BroadcasterUserID+event.UserID]; ok {
			usersForIncrementUsedEmotes[event.BroadcasterUserID+event.UserID].cost += event.Reward.Cost
		} else {
			usersForIncrementUsedEmotes[event.BroadcasterUserID+event.UserID] = &userForIncrementUsedEmotes{
				userId:    event.UserID,
				channelId: event.BroadcasterUserID,
				cost:      event.Reward.Cost,
			}
		}

		// youtube song requests

		go func() {
			e := c.handleYoutubeSongRequests(&event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("e", err))
			}
		}()

		go func() {
			e := c.handleAlerts(&event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("e", err))
			}
		}()

		go func() {
			e := c.handleRewardsSevenTvEmote(&event)
			if e != nil {
				c.logger.Error(e.Error(), slog.Any("err", e))
			}
		}()

		go func() {
			if redemptionAddErr := c.twirBus.RedemptionAdd.Publish(
				twitch.ActivatedRedemption{
					ID:                   event.ID,
					BroadcasterUserID:    event.BroadcasterUserID,
					BroadcasterUserLogin: event.BroadcasterUserLogin,
					BroadcasterUserName:  event.BroadcasterUserName,
					UserID:               event.UserID,
					UserLogin:            event.UserLogin,
					UserName:             event.UserName,
					UserInput:            event.UserInput,
					Status:               event.Status,
					RedeemedAt:           time.Now(),
					Reward: twitch.ActivatedRedemptionReward{
						ID:     event.Reward.ID,
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
				userForIncrement.userId,
				userForIncrement.channelId,
				userForIncrement.cost,
			)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}

		if len(itemsForHistoryCreate) > 0 {
			err := c.redemptionsHistoryRepository.CreateMany(ctx, itemsForHistoryCreate)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}

		if len(itemsForEventsCreate) > 0 {
			err := c.eventsListRepository.CreateMany(ctx, itemsForEventsCreate)
			if err != nil {
				c.logger.Error(err.Error(), slog.Any("err", err))
			}
		}
	}()
}

func (c *Handler) handleChannelPointsRewardRedemptionAdd(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) {
	c.redemptionsBatcher.Add(*event)
}

func (c *Handler) handleChannelPointsRewardRedemptionUpdate(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPointsRewardRedemptionUpdate,
) {
	if event.Status != "CANCELED" {
		return
	}

	userStats := &model.UsersStats{}
	err := c.gorm.Where(`"userId" = ?`, event.UserID).Find(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
	if userStats.ID == "" {
		return
	}
	userStats.UsedChannelPoints -= int64(event.Reward.Cost)
	err = c.gorm.Save(userStats).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}
}

func (c *Handler) countUserChannelPoints(userId, channelId string, count int) error {
	user := &model.Users{}
	err := c.gorm.
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
		err = c.gorm.Save(user.Stats).Error
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
		err = c.gorm.Create(user.Stats).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Handler) handleYoutubeSongRequests(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	if event.UserInput == "" {
		return nil
	}

	entity, err := c.channelSongRequestsSettingsCache.Get(
		context.Background(),
		event.BroadcasterUserID,
	)
	if err != nil {
		return err
	}
	if entity.ID == "" || !entity.Enabled || event.Reward.ID != entity.ChannelPointsRewardID.String {
		return nil
	}

	var foundCommand *model.ChannelsCommands
	commands, err := c.commandsCache.Get(context.Background(), event.BroadcasterUserID)
	if err != nil {
		return err
	}

	for _, command := range commands {
		if command.DefaultName.String == "sr" && command.Enabled {
			foundCommand = &command
			break
		}
	}

	if foundCommand == nil {
		return nil
	}

	res, err := c.twirBus.Parser.GetCommandResponse.Request(
		context.Background(), twitch.TwitchChatMessage{
			BroadcasterUserId:    event.BroadcasterUserID,
			BroadcasterUserName:  event.BroadcasterUserName,
			BroadcasterUserLogin: event.BroadcasterUserLogin,
			ChatterUserId:        event.UserID,
			ChatterUserName:      event.UserName,
			ChatterUserLogin:     event.UserLogin,
			MessageId:            event.ID,
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
			bots.SendMessageRequest{
				ChannelId:      event.BroadcasterUserID,
				ChannelName:    &event.BroadcasterUserLogin,
				Message:        fmt.Sprintf("@%s %s", event.UserLogin, response),
				SkipRateLimits: true,
			},
		)
	}

	return nil
}

func (c *Handler) handleAlerts(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	alerts, err := c.alertsCache.Get(context.Background(), event.BroadcasterUserID)
	if err != nil {
		return err
	}

	var foundAlertId uuid.UUID
	for _, alert := range alerts {
		if slices.Contains(alert.RewardIDS, event.Reward.ID) {
			foundAlertId = alert.ID
			break
		}
	}
	if foundAlertId == uuid.Nil {
		return nil
	}

	_, err = c.websocketsGrpc.TriggerAlert(
		context.TODO(),
		&websockets.TriggerAlertRequest{
			ChannelId: event.BroadcasterUserID,
			AlertId:   foundAlertId.String(),
		},
	)

	return err
}
