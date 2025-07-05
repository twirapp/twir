package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func DudesOverlaySettingsEntityToGql(e entity.DudesOverlaySettings) gqlmodel.DudesOverlaySettings {
	return gqlmodel.DudesOverlaySettings{
		ID: e.ID,
		DudeSettings: &gqlmodel.DudesDudeSettings{
			Color:          e.DudeSettings.Color,
			EyesColor:      e.DudeSettings.EyesColor,
			CosmeticsColor: e.DudeSettings.CosmeticsColor,
			MaxLifeTime:    e.DudeSettings.MaxLifeTime,
			Gravity:        e.DudeSettings.Gravity,
			Scale:          e.DudeSettings.Scale,
			SoundsEnabled:  e.DudeSettings.SoundsEnabled,
			SoundsVolume:   e.DudeSettings.SoundsVolume,
			VisibleName:    e.DudeSettings.VisibleName,
			GrowTime:       e.DudeSettings.GrowTime,
			GrowMaxScale:   e.DudeSettings.GrowMaxScale,
			MaxOnScreen:    e.DudeSettings.MaxOnScreen,
			DefaultSprite:  e.DudeSettings.DefaultSprite,
		},
		MessageBoxSettings: &gqlmodel.DudesMessageBoxSettings{
			Enabled:      e.MessageBoxSettings.Enabled,
			BorderRadius: e.MessageBoxSettings.BorderRadius,
			BoxColor:     e.MessageBoxSettings.BoxColor,
			FontFamily:   e.MessageBoxSettings.FontFamily,
			FontSize:     e.MessageBoxSettings.FontSize,
			Padding:      e.MessageBoxSettings.Padding,
			ShowTime:     e.MessageBoxSettings.ShowTime,
			Fill:         e.MessageBoxSettings.Fill,
		},
		NameBoxSettings: &gqlmodel.DudesNameBoxSettings{
			FontFamily:         e.NameBoxSettings.FontFamily,
			FontSize:           e.NameBoxSettings.FontSize,
			Fill:               e.NameBoxSettings.Fill,
			LineJoin:           e.NameBoxSettings.LineJoin,
			StrokeThickness:    e.NameBoxSettings.StrokeThickness,
			Stroke:             e.NameBoxSettings.Stroke,
			FillGradientStops:  e.NameBoxSettings.FillGradientStops,
			FillGradientType:   e.NameBoxSettings.FillGradientType,
			FontStyle:          e.NameBoxSettings.FontStyle,
			FontVariant:        e.NameBoxSettings.FontVariant,
			FontWeight:         e.NameBoxSettings.FontWeight,
			DropShadow:         e.NameBoxSettings.DropShadow,
			DropShadowAlpha:    e.NameBoxSettings.DropShadowAlpha,
			DropShadowAngle:    e.NameBoxSettings.DropShadowAngle,
			DropShadowBlur:     e.NameBoxSettings.DropShadowBlur,
			DropShadowDistance: e.NameBoxSettings.DropShadowDistance,
			DropShadowColor:    e.NameBoxSettings.DropShadowColor,
		},
		IgnoreSettings: &gqlmodel.DudesIgnoreSettings{
			IgnoreCommands: e.IgnoreSettings.IgnoreCommands,
			IgnoreUsers:    e.IgnoreSettings.IgnoreUsers,
			Users:          e.IgnoreSettings.Users,
		},
		SpitterEmoteSettings: &gqlmodel.DudesSpitterEmoteSettings{
			Enabled: e.SpitterEmoteSettings.Enabled,
		},
	}
}

func DudesOverlaySettingsInputToServiceCreateInput(
	input gqlmodel.DudesOverlaySettingsInput,
	channelID string,
) entity.DudesOverlaySettings {
	return entity.DudesOverlaySettings{
		DudeSettings: entity.DudesDudeSettings{
			Color:          input.DudeSettings.Color,
			EyesColor:      input.DudeSettings.EyesColor,
			CosmeticsColor: input.DudeSettings.CosmeticsColor,
			MaxLifeTime:    input.DudeSettings.MaxLifeTime,
			Gravity:        input.DudeSettings.Gravity,
			Scale:          input.DudeSettings.Scale,
			SoundsEnabled:  input.DudeSettings.SoundsEnabled,
			SoundsVolume:   input.DudeSettings.SoundsVolume,
			VisibleName:    input.DudeSettings.VisibleName,
			GrowTime:       input.DudeSettings.GrowTime,
			GrowMaxScale:   input.DudeSettings.GrowMaxScale,
			MaxOnScreen:    input.DudeSettings.MaxOnScreen,
			DefaultSprite:  input.DudeSettings.DefaultSprite,
		},
		MessageBoxSettings: entity.DudesMessageBoxSettings{
			Enabled:      input.MessageBoxSettings.Enabled,
			BorderRadius: input.MessageBoxSettings.BorderRadius,
			BoxColor:     input.MessageBoxSettings.BoxColor,
			FontFamily:   input.MessageBoxSettings.FontFamily,
			FontSize:     input.MessageBoxSettings.FontSize,
			Padding:      input.MessageBoxSettings.Padding,
			ShowTime:     input.MessageBoxSettings.ShowTime,
			Fill:         input.MessageBoxSettings.Fill,
		},
		NameBoxSettings: entity.DudesNameBoxSettings{
			FontFamily:         input.NameBoxSettings.FontFamily,
			FontSize:           input.NameBoxSettings.FontSize,
			Fill:               input.NameBoxSettings.Fill,
			LineJoin:           input.NameBoxSettings.LineJoin,
			StrokeThickness:    input.NameBoxSettings.StrokeThickness,
			Stroke:             input.NameBoxSettings.Stroke,
			FillGradientStops:  input.NameBoxSettings.FillGradientStops,
			FillGradientType:   input.NameBoxSettings.FillGradientType,
			FontStyle:          input.NameBoxSettings.FontStyle,
			FontVariant:        input.NameBoxSettings.FontVariant,
			FontWeight:         input.NameBoxSettings.FontWeight,
			DropShadow:         input.NameBoxSettings.DropShadow,
			DropShadowAlpha:    input.NameBoxSettings.DropShadowAlpha,
			DropShadowAngle:    input.NameBoxSettings.DropShadowAngle,
			DropShadowBlur:     input.NameBoxSettings.DropShadowBlur,
			DropShadowDistance: input.NameBoxSettings.DropShadowDistance,
			DropShadowColor:    input.NameBoxSettings.DropShadowColor,
		},
		IgnoreSettings: entity.DudesIgnoreSettings{
			IgnoreCommands: input.IgnoreSettings.IgnoreCommands,
			IgnoreUsers:    input.IgnoreSettings.IgnoreUsers,
			Users:          input.IgnoreSettings.Users,
		},
		SpitterEmoteSettings: entity.DudesSpitterEmoteSettings{
			Enabled: input.SpitterEmoteSettings.Enabled,
		},
	}
}

func DudesOverlaySettingsUpdateInputToServiceUpdateInput(input gqlmodel.DudesOverlaySettingsInput) entity.DudesOverlaySettings {
	result := entity.DudesOverlaySettings{}

	if input.DudeSettings != nil {
		result.DudeSettings = entity.DudesDudeSettings{
			Color:          input.DudeSettings.Color,
			EyesColor:      input.DudeSettings.EyesColor,
			CosmeticsColor: input.DudeSettings.CosmeticsColor,
			MaxLifeTime:    input.DudeSettings.MaxLifeTime,
			Gravity:        input.DudeSettings.Gravity,
			Scale:          input.DudeSettings.Scale,
			SoundsEnabled:  input.DudeSettings.SoundsEnabled,
			SoundsVolume:   input.DudeSettings.SoundsVolume,
			VisibleName:    input.DudeSettings.VisibleName,
			GrowTime:       input.DudeSettings.GrowTime,
			GrowMaxScale:   input.DudeSettings.GrowMaxScale,
			MaxOnScreen:    input.DudeSettings.MaxOnScreen,
			DefaultSprite:  input.DudeSettings.DefaultSprite,
		}
	}

	if input.MessageBoxSettings != nil {
		result.MessageBoxSettings = entity.DudesMessageBoxSettings{
			Enabled:      input.MessageBoxSettings.Enabled,
			BorderRadius: input.MessageBoxSettings.BorderRadius,
			BoxColor:     input.MessageBoxSettings.BoxColor,
			FontFamily:   input.MessageBoxSettings.FontFamily,
			FontSize:     input.MessageBoxSettings.FontSize,
			Padding:      input.MessageBoxSettings.Padding,
			ShowTime:     input.MessageBoxSettings.ShowTime,
			Fill:         input.MessageBoxSettings.Fill,
		}
	}

	if input.NameBoxSettings != nil {
		result.NameBoxSettings = entity.DudesNameBoxSettings{
			FontFamily:         input.NameBoxSettings.FontFamily,
			FontSize:           input.NameBoxSettings.FontSize,
			Fill:               input.NameBoxSettings.Fill,
			LineJoin:           input.NameBoxSettings.LineJoin,
			StrokeThickness:    input.NameBoxSettings.StrokeThickness,
			Stroke:             input.NameBoxSettings.Stroke,
			FillGradientStops:  input.NameBoxSettings.FillGradientStops,
			FillGradientType:   input.NameBoxSettings.FillGradientType,
			FontStyle:          input.NameBoxSettings.FontStyle,
			FontVariant:        input.NameBoxSettings.FontVariant,
			FontWeight:         input.NameBoxSettings.FontWeight,
			DropShadow:         input.NameBoxSettings.DropShadow,
			DropShadowAlpha:    input.NameBoxSettings.DropShadowAlpha,
			DropShadowAngle:    input.NameBoxSettings.DropShadowAngle,
			DropShadowBlur:     input.NameBoxSettings.DropShadowBlur,
			DropShadowDistance: input.NameBoxSettings.DropShadowDistance,
			DropShadowColor:    input.NameBoxSettings.DropShadowColor,
		}
	}

	if input.IgnoreSettings != nil {
		result.IgnoreSettings = entity.DudesIgnoreSettings{
			IgnoreCommands: input.IgnoreSettings.IgnoreCommands,
			IgnoreUsers:    input.IgnoreSettings.IgnoreUsers,
			Users:          input.IgnoreSettings.Users,
		}
	}

	if input.SpitterEmoteSettings != nil {
		result.SpitterEmoteSettings = entity.DudesSpitterEmoteSettings{
			Enabled: input.SpitterEmoteSettings.Enabled,
		}
	}

	return result
}
