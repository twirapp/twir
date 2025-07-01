package kappagen

import (
	"errors"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/websockets/internal/protoutils"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/events"
	"github.com/twirapp/twir/libs/api/messages/overlays_kappagen"
)

func (c *Kappagen) SendSettings(userId string) error {
	entity := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "kappagen_overlay").
		Find(entity).
		Error
	if err != nil {
		return err
	}

	if entity.ID == "" {
		return nil
	}

	twitchClient, err := twitch.NewUserClient(userId, c.config, c.twirBus)
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

	parsedEntitySettings := model.KappagenOverlaySettings{}
	err = json.Unmarshal(entity.Settings, &parsedEntitySettings)
	if err != nil {
		return err
	}

	kappagenSettings := overlays_kappagen.Settings{
		Emotes: &overlays_kappagen.Settings_Emotes{
			Time:           parsedEntitySettings.Emotes.Time,
			Max:            parsedEntitySettings.Emotes.Max,
			Queue:          parsedEntitySettings.Emotes.Queue,
			FfzEnabled:     parsedEntitySettings.Emotes.FfzEnabled,
			BttvEnabled:    parsedEntitySettings.Emotes.BttvEnabled,
			SevenTvEnabled: parsedEntitySettings.Emotes.SevenTvEnabled,
			EmojiStyle:     overlays_kappagen.EmojiStyle(parsedEntitySettings.Emotes.EmojiStyle),
		},
		Size: &overlays_kappagen.Settings_Size{
			RatioNormal: parsedEntitySettings.Size.RatioNormal,
			RatioSmall:  parsedEntitySettings.Size.RatioSmall,
			Min:         parsedEntitySettings.Size.Min,
			Max:         parsedEntitySettings.Size.Max,
		},
		Cube: &overlays_kappagen.Settings_Cube{
			Speed: parsedEntitySettings.Cube.Speed,
		},
		Animation: &overlays_kappagen.Settings_Animation{
			FadeIn:  parsedEntitySettings.Animation.FadeIn,
			FadeOut: parsedEntitySettings.Animation.FadeOut,
			ZoomIn:  parsedEntitySettings.Animation.ZoomIn,
			ZoomOut: parsedEntitySettings.Animation.ZoomOut,
		},
		Animations: lo.Map(
			parsedEntitySettings.Animations, func(
				item model.KappagenOverlaySettingsAnimationSettings,
				_ int,
			) *overlays_kappagen.Settings_AnimationSettings {
				var prefs *overlays_kappagen.Settings_AnimationSettings_Prefs
				if item.Prefs != nil {
					prefs = &overlays_kappagen.Settings_AnimationSettings_Prefs{
						Size:    item.Prefs.Size,
						Center:  item.Prefs.Center,
						Speed:   item.Prefs.Speed,
						Faces:   item.Prefs.Faces,
						Message: item.Prefs.Message,
						Time:    item.Prefs.Time,
					}
				}

				return &overlays_kappagen.Settings_AnimationSettings{
					Style:   item.Style,
					Prefs:   prefs,
					Count:   item.Count,
					Enabled: item.Enabled,
				}
			},
		),
		EnableRave: parsedEntitySettings.EnableRave,
		Events: lo.Map(
			parsedEntitySettings.Events,
			func(item model.KappagenOverlaySettingsEvent, _ int) *overlays_kappagen.Settings_Event {
				return &overlays_kappagen.Settings_Event{
					Event:          events.TwirEventType(item.Event),
					DisabledStyles: item.DisabledStyles,
					Enabled:        item.Enabled,
				}
			},
		),
		EnableSpawn:    parsedEntitySettings.EnableSpawn,
		ExcludedEmotes: parsedEntitySettings.ExcludedEmotes,
	}

	d, err := protoutils.CreateJsonWithProto(
		&kappagenSettings,
		map[string]any{
			"channelId":   user.ID,
			"channelName": user.Login,
		},
	)
	if err != nil {
		return err
	}

	return c.SendEvent(
		userId,
		"settings",
		d,
	)
}
