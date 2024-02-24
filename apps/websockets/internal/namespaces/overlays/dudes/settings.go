package dudes

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/websockets/internal/protoutils"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/overlays_dudes"
)

type settings struct {
	model.ChannelsOverlaysDudes
	ChannelID          string `json:"channelId"`
	ChannelName        string `json:"channelName"`
	ChannelDisplayName string `json:"channelDisplayName"`
}

func (c *Dudes) SendSettings(userId string, overlayId string) error {
	entity := model.ChannelsOverlaysDudes{}

	query := c.gorm.Where("channel_id = ?", userId)

	if overlayId != "" {
		query = query.Where("id = ?", overlayId)
	}

	err := query.First(&entity).Error
	if err != nil {
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

	dudesGrpcSettings := overlays_dudes.Settings{
		Id: lo.ToPtr(entity.ID.String()),
		DudeSettings: &overlays_dudes.DudeSettings{
			Color:         entity.DudeColor,
			MaxLifeTime:   entity.DudeMaxLifeTime,
			Gravity:       entity.DudeGravity,
			Scale:         entity.DudeScale,
			SoundsEnabled: entity.DudeSoundsEnabled,
			SoundsVolume:  entity.DudeSoundsVolume,
			VisibleName:   entity.DudeVisibleName,
			GrowTime:      entity.DudeGrowTime,
			GrowMaxScale:  entity.DudeGrowMaxScale,
		},
		MessageBoxSettings: &overlays_dudes.MessageBoxSettings{
			Enabled:      entity.MessageBoxEnabled,
			BorderRadius: entity.MessageBoxBorderRadius,
			BoxColor:     entity.MessageBoxBoxColor,
			FontFamily:   entity.MessageBoxFontFamily,
			FontSize:     entity.MessageBoxFontSize,
			Padding:      entity.MessageBoxPadding,
			ShowTime:     entity.MessageBoxShowTime,
			Fill:         entity.MessageBoxFill,
		},
		NameBoxSettings: &overlays_dudes.NameBoxSettings{
			FontFamily:         entity.NameBoxFontFamily,
			FontSize:           entity.NameBoxFontSize,
			Fill:               entity.NameBoxFill,
			LineJoin:           entity.NameBoxLineJoin,
			StrokeThickness:    entity.NameBoxStrokeThickness,
			Stroke:             entity.NameBoxStroke,
			FillGradientStops:  entity.NameBoxFillGradientStops,
			FillGradientType:   entity.NameBoxFillGradientType,
			FontStyle:          entity.NameBoxFontStyle,
			FontVariant:        entity.NameBoxFontVariant,
			FontWeight:         entity.NameBoxFontWeight,
			DropShadow:         entity.NameBoxDropShadow,
			DropShadowAlpha:    entity.NameBoxDropShadowAlpha,
			DropShadowAngle:    entity.NameBoxDropShadowAngle,
			DropShadowBlur:     entity.NameBoxDropShadowBlur,
			DropShadowDistance: entity.NameBoxDropShadowDistance,
			DropShadowColor:    entity.NameBoxDropShadowColor,
		},
		IgnoreSettings: &overlays_dudes.IgnoreSettings{
			IgnoreCommands: entity.IgnoreCommands,
			IgnoreUsers:    entity.IgnoreUsers,
			Users:          entity.IgnoredUsers,
		},
		SpitterEmoteSettings: &overlays_dudes.SpitterEmoteSettings{
			Enabled: entity.SpitterEmoteEnabled,
		},
	}

	data, err := protoutils.CreateJsonWithProto(
		&dudesGrpcSettings, map[string]any{
			"channelId":          user.ID,
			"channelName":        user.Login,
			"channelDisplayName": user.DisplayName,
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
