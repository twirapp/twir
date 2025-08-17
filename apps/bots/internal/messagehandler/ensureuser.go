package messagehandler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/users"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstats "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
)

func (c *MessageHandler) ensureUser(
	ctx context.Context,
	msg handleMessage,
) (*deprecatedgormmodel.Users, error) {
	handlerCtx, handlerSpan := messageHandlerTracer.Start(ctx, "ensureUser")
	defer handlerSpan.End()

	ctx = handlerCtx

	fetchedUser, err := c.usersRepository.GetByID(ctx, msg.ChatterUserId)
	if err != nil && !errors.Is(err, usermodel.ErrNotFound) {
		return nil, err
	}

	var user *usermodel.User
	var userStats *usersstatsmodel.UserStat

	usedEmotesInMessage := 0
	for _, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
		usedEmotesInMessage += count
	}

	badges := createUserBadges(msg.Badges)
	isMod := lo.Contains(badges, "MODERATOR")
	isSubscriber := lo.Contains(badges, "SUBSCRIBER")
	isVip := lo.Contains(badges, "VIP")

	if fetchedUser == usermodel.Nil {
		err := c.trmManager.Do(
			ctx, func(trCtx context.Context) error {
				createdUser, err := c.usersRepository.Create(
					trCtx, users.CreateInput{
						ID:     msg.ChatterUserId,
						ApiKey: lo.ToPtr(uuid.NewString()),
					},
				)
				if err != nil {
					return err
				}

				user = &createdUser

				messages := 0
				if msg.EnrichedData.ChannelStream != nil {
					messages = 1
				}

				stats, err := c.usersstatsRepository.Create(
					ctx, usersstats.CreateInput{
						UserID:       msg.ChatterUserId,
						ChannelID:    msg.BroadcasterUserId,
						Messages:     int32(messages),
						Emotes:       usedEmotesInMessage,
						IsMod:        isMod,
						IsVip:        isVip,
						IsSubscriber: isSubscriber,
					},
				)

				userStats = stats

				return nil
			},
		)

		if err != nil {
			return nil, err
		}

		return &deprecatedgormmodel.Users{
			ID:         msg.ChatterUserId,
			TokenID:    sql.NullString{},
			IsBotAdmin: user.IsBotAdmin,
			ApiKey:     user.ApiKey,
			Token:      nil,
			Stats: &deprecatedgormmodel.UsersStats{
				ID:                userStats.ID.String(),
				UserID:            msg.ChatterUserId,
				ChannelID:         msg.BroadcasterUserId,
				Messages:          userStats.Messages,
				Watched:           userStats.Watched,
				UsedChannelPoints: userStats.UsedChannelPoints,
				IsMod:             userStats.IsMod,
				IsVip:             userStats.IsVip,
				IsSubscriber:      userStats.IsSubscriber,
				Reputation:        userStats.Reputation,
				Emotes:            userStats.Emotes,
				CreatedAt:         userStats.CreatedAt,
				UpdatedAt:         userStats.UpdatedAt,
			},
			IsBanned:          user.IsBanned,
			CreatedAt:         user.CreatedAt,
			HideOnLandingPage: user.HideOnLandingPage,
		}, nil
	} else {
		user = &fetchedUser

		fetchedStats, err := c.usersstatsRepository.GetByUserAndChannelID(
			ctx,
			msg.ChatterUserId,
			msg.BroadcasterUserId,
		)
		if err != nil && !errors.Is(err, usersstats.ErrNotFound) {
			return nil, err
		}

		if fetchedStats == nil {
			createdStats, err := c.usersstatsRepository.Create(
				ctx, usersstats.CreateInput{
					UserID:    msg.ChatterUserId,
					ChannelID: msg.BroadcasterUserId,
					Messages:  int32(1),
				},
			)
			if err != nil {
				return nil, err
			}
			userStats = createdStats
		} else {
			userStats = fetchedStats

			// if msg.EnrichedData.ChannelStream != nil {
			userStats.Messages += 1
			userStats.Emotes += usedEmotesInMessage
			// }

			if _, err := c.usersstatsRepository.Update(
				ctx,
				msg.ChatterUserId,
				msg.BroadcasterUserId,
				usersstats.UpdateInput{
					NumberFields: usersstats.UpdateNumberFieldsInput{
						usersstats.IncrementInputFieldMessages: {Count: 1, IsIncrement: true},
						usersstats.IncrementInputFieldEmotes:   {Count: usedEmotesInMessage, IsIncrement: true},
					},
					IsMod:        lo.ToPtr(isMod),
					IsVip:        lo.ToPtr(isVip),
					IsSubscriber: lo.ToPtr(isSubscriber),
				},
			); err != nil {
				return nil, err
			}
		}
	}

	return &deprecatedgormmodel.Users{
		ID:         msg.ChatterUserId,
		TokenID:    user.TokenID.NullString,
		IsBotAdmin: fetchedUser.IsBotAdmin,
		ApiKey:     fetchedUser.ApiKey,
		Stats: &deprecatedgormmodel.UsersStats{
			ID:                userStats.ID.String(),
			UserID:            msg.ChatterUserId,
			ChannelID:         msg.BroadcasterUserId,
			Messages:          userStats.Messages,
			Watched:           userStats.Watched,
			UsedChannelPoints: userStats.UsedChannelPoints,
			IsMod:             userStats.IsMod,
			IsVip:             userStats.IsVip,
			IsSubscriber:      userStats.IsSubscriber,
			Reputation:        userStats.Reputation,
			Emotes:            userStats.Emotes,
			CreatedAt:         userStats.CreatedAt,
			UpdatedAt:         userStats.UpdatedAt,
		},
		IsBanned:          fetchedUser.IsBanned,
		CreatedAt:         fetchedUser.CreatedAt,
		HideOnLandingPage: fetchedUser.HideOnLandingPage,
	}, nil
}
