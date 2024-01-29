package chat

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/websockets/internal/protoutils"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/overlays_chat"
	"gorm.io/gorm"
)

func (c *Chat) SendSettings(userId string, overlayId string) error {
	entity := model.ChatOverlaySettings{}

	query := c.gorm.Where(`"channel_id" = ?`, userId)

	if overlayId != "" {
		query = query.Where("id = ?", overlayId)
	}

	err := query.Order("created_at asc").First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return err
	}

	twitchClient, err := twitch.NewUserClient(userId, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	usersReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return err
	}

	if len(usersReq.Data.Users) == 0 {
		return errors.New("cannot get user")
	}

	user := usersReq.Data.Users[0]

	channelBadgesReq, err := twitchClient.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: userId,
		},
	)
	if err != nil {
		return err
	}

	globalBadgesReq, err := twitchClient.GetGlobalChatBadges()
	if err != nil {
		return err
	}

	overlaySettings := overlays_chat.Settings{
		Id:                  lo.ToPtr(entity.ID.String()),
		MessageHideTimeout:  entity.MessageHideTimeout,
		MessageShowDelay:    entity.MessageShowDelay,
		Preset:              entity.Preset,
		FontSize:            entity.FontSize,
		HideCommands:        entity.HideCommands,
		HideBots:            entity.HideBots,
		FontFamily:          entity.FontFamily,
		ShowBadges:          entity.ShowBadges,
		ShowAnnounceBadge:   entity.ShowAnnounceBadge,
		TextShadowColor:     entity.TextShadowColor,
		TextShadowSize:      entity.TextShadowSize,
		ChatBackgroundColor: entity.ChatBackgroundColor,
		Direction:           entity.Direction,
		FontWeight:          entity.FontWeight,
		FontStyle:           entity.FontStyle,
		PaddingContainer:    entity.PaddingContainer,
	}

	globalBadges := map[string]helix.ChatBadge{}
	for _, badge := range globalBadgesReq.Data.Badges {
		globalBadges[badge.SetID] = badge
	}

	channelBadges := map[string]helix.ChatBadge{}
	for _, badge := range channelBadgesReq.Data.Badges {
		channelBadges[badge.SetID] = badge
	}

	data, err := protoutils.CreateJsonWithProto(
		&overlaySettings, map[string]any{
			"channelId":          user.ID,
			"channelName":        user.Login,
			"channelDisplayName": user.DisplayName,
			"globalBadges":       globalBadges,
			"channelBadges":      channelBadges,
		},
	)
	if err != nil {
		return err
	}

	return c.SendEvent(
		userId,
		"settings",
		data,
	)
}
