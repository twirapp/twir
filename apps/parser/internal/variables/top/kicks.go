package top

import (
	"context"
	"strconv"
	"strings"

	"github.com/scorfly/gokick"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	sharedvars "github.com/twirapp/twir/apps/parser/internal/variables/shared"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

var KicksLifetime = kickLeaderboardVariable("top.kicks.lifetime", func(data kickLeaderboardData) []kickLeaderboardEntry {
	return data.Lifetime
})

var KicksMonth = kickLeaderboardVariable("top.kicks.month", func(data kickLeaderboardData) []kickLeaderboardEntry {
	return data.Month
})

var KicksWeek = kickLeaderboardVariable("top.kicks.week", func(data kickLeaderboardData) []kickLeaderboardEntry {
	return data.Week
})

type kickLeaderboardEntry struct {
	Username     string
	GiftedAmount int
}

type kickLeaderboardData struct {
	Lifetime []kickLeaderboardEntry
	Month    []kickLeaderboardEntry
	Week     []kickLeaderboardEntry
}

var kickLeaderboardRequester = func(ctx context.Context, parseCtx *types.VariableParseContext, top int) (*kickLeaderboardData, error) {
	resp, err := sharedvars.GetKickKicksLeaderboard(ctx, parseCtx, top)
	if err != nil {
		return nil, err
	}

	return &kickLeaderboardData{
		Lifetime: lo.Map(resp.Lifetime, func(item gokick.KicksLeaderboardEntry, _ int) kickLeaderboardEntry {
			return kickLeaderboardEntry{Username: item.Username, GiftedAmount: item.GiftedAmount}
		}),
		Month: lo.Map(resp.Month, func(item gokick.KicksLeaderboardEntry, _ int) kickLeaderboardEntry {
			return kickLeaderboardEntry{Username: item.Username, GiftedAmount: item.GiftedAmount}
		}),
		Week: lo.Map(resp.Week, func(item gokick.KicksLeaderboardEntry, _ int) kickLeaderboardEntry {
			return kickLeaderboardEntry{Username: item.Username, GiftedAmount: item.GiftedAmount}
		}),
	}, nil
}

func SetKickLeaderboardRequesterForTests(requester func(ctx context.Context, parseCtx *types.VariableParseContext, top int) (*kickLeaderboardData, error)) func() {
	old := kickLeaderboardRequester
	kickLeaderboardRequester = requester
	return func() {
		kickLeaderboardRequester = old
	}
}

func kickLeaderboardVariable(name string, selector func(kickLeaderboardData) []kickLeaderboardEntry) *types.Variable {
	return &types.Variable{
		Name:                name,
		Description:         lo.ToPtr("Kick gifted kicks leaderboard"),
		CanBeUsedInRegistry: true,
		Handler: func(
			ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := &types.VariableHandlerResult{}
			if parseCtx.Platform != platformentity.PlatformKick {
				result.Result = "not supported on this platform"
				return result, nil
			}

			limit := 10
			if variableData.Params != nil {
				if newLimit, err := strconv.Atoi(*variableData.Params); err == nil && newLimit > 0 && newLimit <= 100 {
					limit = newLimit
				}
			}

			data, err := kickLeaderboardRequester(ctx, parseCtx, limit)
			if err != nil {
				return nil, err
			}

			entries := selector(*data)
			if len(entries) == 0 {
				return result, nil
			}

			mapped := lo.Map(entries, func(item kickLeaderboardEntry, _ int) string {
				return item.Username + " × " + strconv.Itoa(item.GiftedAmount)
			})

			result.Result = strings.Join(mapped, " · ")
			return result, nil
		},
	}
}
