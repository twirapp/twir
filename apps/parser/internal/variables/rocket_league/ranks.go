package rocket_league

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
)

var Rank2v2 = &types.Variable{
	Name: "rocketleague.profile.2v2rank",
	Description: lo.ToPtr(
		`You rank in 2v2 ranked, example: Grand Champion II Division IV"`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Doubles 2v2"
		})
		if !ok {
			return nil, nil
		}
		result.Result = rank.Rank

		return &result, nil
	},
}

var Rating2v2 = &types.Variable{
	Name: "rocketleague.profile.2v2rating",
	Description: lo.ToPtr(
		`You rating in 2v2 ranked, example: 1500`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Doubles 2v2"
		})
		if !ok {
			return nil, nil
		}
		result.Result = fmt.Sprintf("%v", rank.Rating)

		return &result, nil
	},
}

var Rank1v1 = &types.Variable{
	Name: "rocketleague.profile.1v1rank",
	Description: lo.ToPtr(
		`You rank in 1v1 ranked, example: Grand Champion II Division IV"`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Duel 1v1"
		})
		if !ok {
			return nil, nil
		}
		result.Result = rank.Rank

		return &result, nil
	},
}

var Rating1v1 = &types.Variable{
	Name: "rocketleague.profile.1v1rating",
	Description: lo.ToPtr(
		`You rating in 1v1 ranked, example: 1500`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Duel 1v1"
		})
		if !ok {
			return nil, nil
		}
		result.Result = fmt.Sprintf("%v", rank.Rating)

		return &result, nil
	},
}

var Rank3v3 = &types.Variable{
	Name: "rocketleague.profile.3v3rank",
	Description: lo.ToPtr(
		`You rank in 3v3 ranked, example: Grand Champion II Division IV"`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Standard 3v3"
		})
		if !ok {
			return nil, nil
		}
		result.Result = rank.Rank

		return &result, nil
	},
}

var Rating3v3 = &types.Variable{
	Name: "rocketleague.profile.1v1rating",
	Description: lo.ToPtr(
		`You rating in 3v3 ranked, example: 1500`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Ranked Standard 3v3"
		})
		if !ok {
			return nil, nil
		}
		result.Result = fmt.Sprintf("%v", rank.Rating)

		return &result, nil
	},
}

var RankRumble = &types.Variable{
	Name: "rocketleague.profile.rumbleRank",
	Description: lo.ToPtr(
		`You rank in rumble, example: Grand Champion II Division IV"`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Rumble"
		})
		if !ok {
			return nil, nil
		}
		result.Result = rank.Rank

		return &result, nil
	},
}

var RatingRumble = &types.Variable{
	Name: "rocketleague.profile.rubmleRating",
	Description: lo.ToPtr(
		`You rating in dropshot, example: 1500`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Rumble"
		})
		if !ok {
			return nil, nil
		}
		result.Result = fmt.Sprintf("%v", rank.Rating)

		return &result, nil
	},
}

var RankDropshot = &types.Variable{
	Name: "rocketleague.profile.dropshotRank",
	Description: lo.ToPtr(
		`You rank in dropshot, example: Grand Champion II Division IV"`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Dropshot"
		})
		if !ok {
			return nil, nil
		}
		result.Result = rank.Rank

		return &result, nil
	},
}

var RatingDropshot = &types.Variable{
	Name: "rocketleague.profile.dropshotRating",
	Description: lo.ToPtr(
		`You rating in rumble, example: 1500`,
	),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}
		stats := parseCtx.Cacher.GetRocketLeagueUserStats(ctx)
		if stats == nil {
			return nil, nil
		}
		rank, ok := lo.Find(stats.Rankings, func(rank types.Ranking) bool {
			return rank.Playlist == "Dropshot"
		})
		if !ok {
			return nil, nil
		}
		result.Result = fmt.Sprintf("%v", rank.Rating)

		return &result, nil
	},
}
