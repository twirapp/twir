package valorant_matches

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"
)

//for (const match of data) {
//if (!match.players?.all_players) {
//continue;
//}
//const player = match.players.all_players.find(p => p.name === 'IW 7ssk7');
//const teamName = player.team.toLowerCase();
//const team = match.teams[teamName];
//const isWin = team.has_won;
//const char = player.character;
//const KDA = `${player.stats.kills}/${player.stats.deaths}/${player.stats.assists}`;
//const matchResultString = isWin ? 'W' : 'L';
//
//result.push(`${matchResultString}(${team.rounds_won}/${team.rounds_lost}) — ${char } ${KDA}`);
//}
//
//return result.join(' | ');

var Trend = types.Variable{
	Name: "valorant.matches.trend",
	Description: lo.ToPtr(
		`Latest 5 matches trend, i.e "W(13/4) — Killjoy 12/4/10 | L(4/13) — Killjoy 4/12/10"`,
	),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		integrations := ctx.GetEnabledIntegrations()
		integration, ok := lo.Find(integrations, func(item model.ChannelsIntegrations) bool {
			return item.Integration.Service == "VALORANT"
		})

		if !ok || integration.Data == nil || integration.Data.UserName == nil {
			return nil, nil
		}

		matches := ctx.GetValorantMatches()
		if len(matches) == 0 {
			return nil, nil
		}

		var trend []string

		for _, match := range matches {
			if len(match.Players.AllPlayers) == 0 {
				continue
			}

			player, ok := lo.Find(match.Players.AllPlayers, func(el variables_cache.ValorantMatchPlayer) bool {
				return fmt.Sprintf("%s#%s", el.Name, el.Tag) == *integration.Data.UserName
			})

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
				fmt.Sprintf(
					"%s(%d/%d) — %s %s",
					matchResultString,
					team.RoundsWon,
					team.RoundsLost,
					char,
					KDA,
				),
			)
		}

		result.Result = strings.Join(trend, " · ")

		return &result, nil
	},
}
