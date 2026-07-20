package dota2

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"gorm.io/gorm"
)

type medalTier string

const (
	medalHerald   medalTier = "herald"
	medalGuardian medalTier = "guardian"
	medalCrusader medalTier = "crusader"
	medalArchon   medalTier = "archon"
	medalLegend   medalTier = "legend"
	medalAncient  medalTier = "ancient"
	medalDivine   medalTier = "divine"
	medalImmortal medalTier = "immortal"
)

func medalForMMR(mmr int) medalTier {
	switch {
	case mmr < 770:
		return medalHerald
	case mmr < 1540:
		return medalGuardian
	case mmr < 2310:
		return medalCrusader
	case mmr < 3080:
		return medalArchon
	case mmr < 3850:
		return medalLegend
	case mmr < 4620:
		return medalAncient
	case mmr < 5420:
		return medalDivine
	default:
		return medalImmortal
	}
}

type winLossOutput struct {
	Record  string
	WinRate string
}

func formatWinLoss(wins, losses int) winLossOutput {
	games := wins + losses
	winRate := 0.0
	if games > 0 {
		winRate = float64(wins) / float64(games) * 100
	}

	return winLossOutput{
		Record:  fmt.Sprintf("%d-%d", wins, losses),
		WinRate: fmt.Sprintf("%.1f%%", winRate),
	}
}

func formatDuration(seconds int) string {
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}

type lastGameOutput struct {
	HeroName string
	KDA      string
	Won      bool
	Duration string
}

func formatLastGame(lastGame *busdota.LastGameInfo) (lastGameOutput, bool) {
	if lastGame == nil {
		return lastGameOutput{}, false
	}

	return lastGameOutput{
		HeroName: lastGame.HeroName,
		KDA:      fmt.Sprintf("%d/%d/%d", lastGame.Kills, lastGame.Deaths, lastGame.Assists),
		Won:      lastGame.Win,
		Duration: formatDuration(lastGame.DurationS),
	}, true
}

func formatWinProbability(probability float64) string {
	return fmt.Sprintf("%.1f%%", probability*100)
}

func joinNotablePlayers(players []string) string {
	return strings.Join(players, ", ")
}

func requireDotaSettings(
	ctx context.Context,
	parseCtx *types.ParseContext,
	isCommandEnabled func(model.ChannelsDotaSettingsCommands) bool,
) (*model.ChannelsDotaSettings, error) {
	settings := &model.ChannelsDotaSettings{}
	err := parseCtx.Services.Gorm.WithContext(ctx).
		Where("channel_id = ?", parseCtx.Channel.DBChannelID).
		First(settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.SettingsNotFound),
				Err:     err,
			}
		}

		return nil, &types.CommandHandlerError{
			Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.GetSettings),
			Err:     fmt.Errorf("get Dota settings: %w", err),
		}
	}

	if !settings.Enabled {
		return nil, &types.CommandHandlerError{
			Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.Disabled),
		}
	}

	if isCommandEnabled != nil && !isCommandEnabled(settings.CommandsSettings) {
		return nil, &types.CommandHandlerError{
			Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.CommandDisabled),
		}
	}

	return settings, nil
}

func getDotaData(
	ctx context.Context,
	parseCtx *types.ParseContext,
) (*busdota.GetDataResponse, error) {
	response, err := parseCtx.Services.Bus.Dota.GetData.Request(
		ctx,
		busdota.GetDataRequest{
			ChannelID:    parseCtx.Channel.DBChannelID,
			TwitchUserID: parseCtx.Channel.ID,
		},
	)
	if err != nil {
		return nil, &types.CommandHandlerError{
			Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.GetData),
			Err:     fmt.Errorf("get Dota data: %w", err),
		}
	}

	if response == nil {
		return nil, &types.CommandHandlerError{
			Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.GetData),
			Err:     errors.New("empty Dota data response"),
		}
	}

	return &response.Data, nil
}

func medalName(ctx context.Context, medal medalTier) string {
	switch medal {
	case medalHerald:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Herald)
	case medalGuardian:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Guardian)
	case medalCrusader:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Crusader)
	case medalArchon:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Archon)
	case medalLegend:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Legend)
	case medalAncient:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Ancient)
	case medalDivine:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Divine)
	default:
		return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Medals.Immortal)
	}
}

func LocalizeDescriptions(ctx context.Context) {
	Mmr.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Mmr),
	)
	MmrSet.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.MmrSet),
	)
	Wl.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Wl),
	)
	Lg.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Lg),
	)
	Gm.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Gm),
	)
	Np.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Np),
	)
	Wp.ChannelsCommands.Description = null.StringFrom(
		i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Descriptions.Wp),
	)
}
