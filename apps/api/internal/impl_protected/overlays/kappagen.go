package overlays

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/events"
	"github.com/satont/twir/libs/grpc/generated/api/overlays_kappagen"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

const kappagenOverlayType = "kappagen_overlay"

func (c *Overlays) kappagenDbToGrpc(s model.KappagenOverlaySettings) *overlays_kappagen.Settings {
	return &overlays_kappagen.Settings{
		Emotes: &overlays_kappagen.Settings_Emotes{
			Time:           s.Emotes.Time,
			Max:            s.Emotes.Max,
			Queue:          s.Emotes.Queue,
			FfzEnabled:     s.Emotes.FfzEnabled,
			BttvEnabled:    s.Emotes.BttvEnabled,
			SevenTvEnabled: s.Emotes.SevenTvEnabled,
			EmojiStyle:     overlays_kappagen.EmojiStyle(s.Emotes.EmojiStyle),
		},
		Size: &overlays_kappagen.Settings_Size{
			RatioNormal: s.Size.RatioNormal,
			RatioSmall:  s.Size.RatioSmall,
			Min:         s.Size.Min,
			Max:         s.Size.Max,
		},
		Cube: &overlays_kappagen.Settings_Cube{
			Speed: s.Cube.Speed,
		},
		Animation: &overlays_kappagen.Settings_Animation{
			FadeIn:  s.Animation.FadeIn,
			FadeOut: s.Animation.FadeOut,
			ZoomIn:  s.Animation.ZoomIn,
			ZoomOut: s.Animation.ZoomOut,
		},
		Animations: lo.Map(
			s.Animations, func(
				v model.KappagenOverlaySettingsAnimationSettings,
				i int,
			) *overlays_kappagen.Settings_AnimationSettings {
				return &overlays_kappagen.Settings_AnimationSettings{
					Style: v.Style,
					Prefs: lo.IfF(
						v.Prefs != nil,
						func() *overlays_kappagen.Settings_AnimationSettings_Prefs {
							return &overlays_kappagen.Settings_AnimationSettings_Prefs{
								Size:    v.Prefs.Size,
								Center:  v.Prefs.Center,
								Speed:   v.Prefs.Speed,
								Faces:   v.Prefs.Faces,
								Message: v.Prefs.Message,
								Time:    v.Prefs.Time,
							}
						},
					).Else(nil),
					Count:   v.Count,
					Enabled: v.Enabled,
				}
			},
		),
		EnableRave: s.EnableRave,
		Events: lo.Map(
			s.Events, func(
				item model.KappagenOverlaySettingsEvent,
				_ int,
			) *overlays_kappagen.Settings_Event {
				return &overlays_kappagen.Settings_Event{
					Event:          events.TwirEventType(item.Event),
					DisabledStyles: item.DisabledStyles,
					Enabled:        item.Enabled,
				}
			},
		),
		EnableSpawn:    s.EnableSpawn,
		ExcludedEmotes: s.ExcludedEmotes,
	}
}

func (c *Overlays) kappagenGrpcToDb(s *overlays_kappagen.Settings) model.KappagenOverlaySettings {
	return model.KappagenOverlaySettings{
		Emotes: model.KappagenOverlaySettingsEmotes{
			Time:           s.Emotes.Time,
			Max:            s.Emotes.Max,
			Queue:          s.Emotes.Queue,
			FfzEnabled:     s.Emotes.FfzEnabled,
			BttvEnabled:    s.Emotes.BttvEnabled,
			SevenTvEnabled: s.Emotes.SevenTvEnabled,
			EmojiStyle:     model.KappagenOverlaySettingsEmotesEmojiStyle(s.Emotes.EmojiStyle),
		},
		Size: model.KappagenOverlaySettingsSize{
			RatioNormal: s.Size.RatioNormal,
			RatioSmall:  s.Size.RatioSmall,
			Min:         s.Size.Min,
			Max:         s.Size.Max,
		},
		Cube: model.KappagenOverlaySettingsCube{
			Speed: s.Cube.Speed,
		},
		Animation: model.KappagenOverlaySettingsAnimation{
			FadeIn:  s.Animation.FadeIn,
			FadeOut: s.Animation.FadeOut,
			ZoomIn:  s.Animation.ZoomIn,
			ZoomOut: s.Animation.ZoomOut,
		},
		Animations: lo.Map(
			s.Animations, func(
				v *overlays_kappagen.Settings_AnimationSettings,
				i int,
			) model.KappagenOverlaySettingsAnimationSettings {
				return model.KappagenOverlaySettingsAnimationSettings{
					Style: v.Style,
					Prefs: lo.IfF(
						v.Prefs != nil,
						func() *model.KappagenOverlaySettingsAnimationSettingsPrefs {
							return &model.KappagenOverlaySettingsAnimationSettingsPrefs{
								Size:    v.Prefs.Size,
								Center:  v.Prefs.Center,
								Speed:   v.Prefs.Speed,
								Faces:   v.Prefs.Faces,
								Message: v.Prefs.Message,
								Time:    v.Prefs.Time,
							}
						},
					).Else(nil),
					Count:   v.Count,
					Enabled: v.Enabled,
				}
			},
		),
		EnableRave: s.EnableRave,
		Events: lo.Map(
			s.Events,
			func(
				item *overlays_kappagen.Settings_Event,
				_ int,
			) model.KappagenOverlaySettingsEvent {
				return model.KappagenOverlaySettingsEvent{
					Event:          int32(item.Event),
					DisabledStyles: item.DisabledStyles,
					Enabled:        item.Enabled,
				}
			},
		),
		EnableSpawn:    s.EnableSpawn,
		ExcludedEmotes: s.ExcludedEmotes,
	}
}

func (c *Overlays) OverlayKappaGenGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*overlays_kappagen.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelModulesSettings{}

	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		First(&entity).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, twirp.NotFoundError("settings not found")
		}
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	parsedSettings := model.KappagenOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	return c.kappagenDbToGrpc(parsedSettings), nil
}

func (c *Overlays) OverlayKappaGenUpdate(
	ctx context.Context,
	req *overlays_kappagen.Settings,
) (*overlays_kappagen.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? and type = ?`, dashboardId, kappagenOverlayType).
		Find(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get settings: %w", err)
	}

	if entity.ID == "" {
		entity.ID = uuid.NewString()
		entity.ChannelId = dashboardId
		entity.Type = kappagenOverlayType
	}

	parsedSettings := c.kappagenGrpcToDb(req)
	settingsJson, err := json.Marshal(parsedSettings)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal settings: %w", err)
	}
	entity.Settings = settingsJson
	if err := c.Db.
		WithContext(ctx).
		Save(&entity).
		Error; err != nil {
		return nil, fmt.Errorf("cannot update settings: %w", err)
	}

	newSettings := model.KappagenOverlaySettings{}
	if err := json.Unmarshal(entity.Settings, &newSettings); err != nil {
		return nil, fmt.Errorf("cannot parse settings: %w", err)
	}

	c.Grpc.Websockets.RefreshKappagenOverlaySettings(
		ctx, &websockets.RefreshKappagenOverlaySettingsRequest{
			ChannelId: dashboardId,
		},
	)

	return c.kappagenDbToGrpc(newSettings), nil
}
