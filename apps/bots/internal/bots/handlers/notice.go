package handlers

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/satont/twir/libs/grpc/generated/events"
)

type OnNoticeOpts struct {
	ChannelID   string
	UserID      string
	Type        string
	Message     string
	RaidViewers string

	SenderUserLogin   string
	SenderDisplayName string
}

func (c *Handlers) OnNotice(opts OnNoticeOpts, rawTags map[string]string) {
	if opts.Type == "raid" {
		viewers := rawTags["msg-param-viewerCount"]
		intViewers, _ := strconv.Atoi(viewers)

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: opts.ChannelID,
				UserID:    opts.UserID,
				Type:      model.ChannelEventListItemTypeRaided,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					RaidedViewersCount:    viewers,
					RaidedFromDisplayName: opts.SenderDisplayName,
					RaidedFromUserName:    opts.SenderUserLogin,
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create raid entity",
				slog.String("channelId", opts.ChannelID),
				slog.String("userId", opts.UserID),
			)
		}
		if _, err := c.eventsGrpc.Raided(
			context.Background(), &events.RaidedMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: opts.ChannelID},
				UserName:        opts.SenderUserLogin,
				UserDisplayName: opts.SenderDisplayName,
				Viewers:         int64(intViewers),
				UserId:          opts.UserID,
			},
		); err != nil {
			c.logger.Error(
				"cannot fire raid event",
				slog.String("channelId", opts.ChannelID),
				slog.String("userId", opts.UserID),
			)
		}
	}

	if rawTags["msg-id"] == "sub" {
		level := getSubPlan(rawTags["msg-param-sub-plan"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: opts.ChannelID,
				Type:      model.ChannelEventListItemTypeSubscribe,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					SubUserName:        rawTags["login"],
					SubUserDisplayName: rawTags["display-name"],
					SubLevel:           level,
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create sub entity",
				slog.String("channelId", opts.ChannelID),
				slog.String("userId", opts.UserID),
			)
		}
		if _, err := c.eventsGrpc.Subscribe(
			context.Background(), &events.SubscribeMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: opts.ChannelID},
				UserName:        rawTags["login"],
				UserDisplayName: rawTags["display-name"],
				Level:           level,
				UserId:          opts.UserID,
			},
		); err != nil {
			c.logger.Error(
				"cannot fire sub event",
				slog.String("channelId", opts.ChannelID),
				slog.String("userId", opts.UserID),
			)
		}
	}

	if rawTags["msg-id"] == "resub" {
		level := getSubPlan(rawTags["msg-param-sub-plan"])
		months, _ := strconv.Atoi(rawTags["msg-param-streak-months"])
		streak, _ := strconv.Atoi(rawTags["msg-param-cumulative-months"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: opts.ChannelID,
				Type:      model.ChannelEventListItemTypeReSubscribe,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					ReSubUserName:        rawTags["login"],
					ReSubUserDisplayName: rawTags["display-name"],
					ReSubLevel:           level,
					ReSubStreak:          rawTags["msg-param-streak-months"],
					ReSubMonths:          rawTags["msg-param-cumulative-months"],
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create resub entity",
				slog.String("channelId", opts.ChannelID),
				slog.String("userLogin", rawTags["login"]),
			)
		}
		if _, err := c.eventsGrpc.ReSubscribe(
			context.Background(), &events.ReSubscribeMessage{
				BaseInfo:        &events.BaseInfo{ChannelId: opts.ChannelID},
				UserName:        rawTags["login"],
				UserDisplayName: rawTags["display-name"],
				Months:          int64(months),
				Streak:          int64(streak),
				IsPrime:         level == "prime",
				Message:         opts.Message,
				Level:           level,
				UserId:          opts.UserID,
			},
		); err != nil {
			c.logger.Error(
				"cannot fire resub entity",
				slog.String("channelId", opts.ChannelID),
				slog.String("userLogin", rawTags["login"]),
				slog.String("userId", opts.UserID),
			)
		}
	}

	if rawTags["msg-id"] == "subgift" {
		level := getSubPlan(rawTags["msg-param-sub-plan"])

		if err := c.db.Create(
			&model.ChannelsEventsListItem{
				ID:        uuid.New().String(),
				ChannelID: opts.ChannelID,
				Type:      model.ChannelEventListItemTypeSubGift,
				CreatedAt: time.Now().UTC(),
				Data: &model.ChannelsEventsListItemData{
					SubGiftLevel:                 level,
					SubGiftUserName:              rawTags["login"],
					SubGiftUserDisplayName:       rawTags["display-name"],
					SubGiftTargetUserName:        rawTags["msg-param-recipient-user-name"],
					SubGiftTargetUserDisplayName: rawTags["msg-param-recipient-display-name"],
				},
			},
		).Error; err != nil {
			c.logger.Error(
				"cannot create subgift entity",
				slog.String("channelId", opts.ChannelID),
				slog.String("userId", opts.UserID),
			)
		}
		c.eventsGrpc.SubGift(
			context.Background(), &events.SubGiftMessage{
				BaseInfo:          &events.BaseInfo{ChannelId: opts.ChannelID},
				SenderUserName:    rawTags["login"],
				SenderDisplayName: rawTags["display-name"],
				TargetUserName:    rawTags["msg-param-recipient-user-name"],
				TargetDisplayName: rawTags["msg-param-recipient-display-name"],
				Level:             level,
				SenderUserId:      opts.UserID,
			},
		)
	}
}

func getSubPlan(plan string) string {
	if plan == "Prime" {
		return "prime"
	}

	parsedPlan, err := strconv.Atoi(plan)
	if err != nil {
		return plan
	}

	return strconv.Itoa(parsedPlan / 1000)
}
