package valorant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var Matches = &types.Variable{
	Name: "valorant.matches.trend",
	Description: lo.ToPtr(
		`Latest 5 matches trend, i.e "W(13/4) — Killjoy 12/4/10 | L(4/13) — Killjoy 4/12/10"`,
	),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		integrations := parseCtx.Cacher.GetEnabledChannelIntegrations(ctx)
		integration, ok := lo.Find(
			integrations, func(item *model.ChannelsIntegrations) bool {
				return item.Integration.Service == "VALORANT"
			},
		)

		if !ok || integration.Data == nil || integration.Data.UserName == nil ||
			integration.Data.ValorantPuuid == nil {
			return &result, nil
		}

		matches := parseCtx.Cacher.GetValorantMatches(ctx)
		if len(matches) == 0 {
			return &result, nil
		}

		if parseCtx.ChannelStream != nil {
			streamMatches := lo.Filter(matches, func(match types.ValorantMatch, _ int) bool {
				t := time.Unix(int64(match.MetaData.GameStart), 0)
				return t.After(parseCtx.ChannelStream.StartedAt.Add(-time.Hour))
			})

			if len(streamMatches) > 0 {
				matches = streamMatches
			}
		}

		var trend []string

		for _, match := range matches {
			if len(match.Players.AllPlayers) == 0 {
				continue
			}

			player, ok := lo.Find(
				match.Players.AllPlayers, func(el types.ValorantMatchPlayer) bool {
					return el.Puuid == *integration.Data.ValorantPuuid
				},
			)

			if !ok {
				continue
			}

			teamName := strings.ToLower(player.Team)
			team := match.Teams[teamName]
			isWin := team.HasWon
			char := player.Character
			KDA := fmt.Sprintf("%d/%d/%d", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists)
			matchResultString := "W"
			if !isWin {
				matchResultString = "L"
			}

			trend = append(
				trend,
				i18n.GetCtx(
					ctx,
					locales.Translations.Variables.Valorant.Info.Matches.
						SetVars(locales.KeysVariablesValorantInfoMatchesVars{MatchResult: matchResultString, RoundsWon: team.RoundsWon, RoundsLost: team.RoundsLost, Char: char, KDA: KDA}),
				),
			)
		}

		result.Result = strings.Join(trend, " · ")

		return &result, nil
	},
}

var MatchesSimple = &types.Variable{
	Name: "valorant.matches.trend.simple",
	Description: lo.ToPtr(
		`Latest 5 matches trend, i.e "5 W / 2 L"`,
	),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		integrations := parseCtx.Cacher.GetEnabledChannelIntegrations(ctx)
		integration, ok := lo.Find(
			integrations, func(item *model.ChannelsIntegrations) bool {
				return item.Integration.Service == "VALORANT"
			},
		)

		if !ok || integration.Data == nil || integration.Data.UserName == nil ||
			integration.Data.ValorantPuuid == nil {
			return &result, nil
		}

		matches := parseCtx.Cacher.GetValorantMatches(ctx)
		if len(matches) == 0 {
			return &result, nil
		}

		if parseCtx.ChannelStream != nil {
			streamMatches := lo.Filter(matches, func(match types.ValorantMatch, _ int) bool {
				t := time.Unix(int64(match.MetaData.GameStart), 0)
				return t.After(parseCtx.ChannelStream.StartedAt.Add(-time.Hour))
			})

			if len(streamMatches) > 0 {
				matches = streamMatches
			}
		}

		var (
			wins  int8
			loses int8
		)

		for _, match := range matches {
			if len(match.Players.AllPlayers) == 0 {
				continue
			}

			player, ok := lo.Find(
				match.Players.AllPlayers, func(el types.ValorantMatchPlayer) bool {
					return el.Puuid == *integration.Data.ValorantPuuid
				},
			)

			if !ok {
				continue
			}

			teamName := strings.ToLower(player.Team)
			team := match.Teams[teamName]
			isWin := team.HasWon

			if isWin {
				wins++
			} else {
				loses++
			}
		}

		result.Result = fmt.Sprintf("%d W / %d L", wins, loses)

		return &result, nil
	},
}
