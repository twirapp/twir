package top

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/parser/internal/types"
	sharedvars "github.com/twirapp/twir/apps/parser/internal/variables/shared"
	parserservices "github.com/twirapp/twir/apps/parser/internal/types/services"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestKickLeaderboardVariables(t *testing.T) {
	restore := SetKickLeaderboardRequesterForTests(func(ctx context.Context, parseCtx *types.VariableParseContext, top int) (*kickLeaderboardData, error) {
		require.Equal(t, 2, top)
		return &kickLeaderboardData{
			Lifetime: []kickLeaderboardEntry{{Username: "foo", GiftedAmount: 100}, {Username: "bar", GiftedAmount: 50}},
			Month:    []kickLeaderboardEntry{{Username: "baz", GiftedAmount: 30}},
			Week:     []kickLeaderboardEntry{{Username: "qux", GiftedAmount: 10}},
		}, nil
	})
	t.Cleanup(restore)

	parseCtx := &types.VariableParseContext{ParseContext: &types.ParseContext{
		Platform: sharedvars.PlatformKick,
		Channel:  &types.ParseContextChannel{ID: "123", Name: "kick-channel"},
		Services: &parserservices.Services{},
	}}
	params := "2"

	lifetime, err := KicksLifetime.Handler(context.Background(), parseCtx, &types.VariableData{Key: "top.kicks.lifetime", Params: &params})
	require.NoError(t, err)
	require.Equal(t, "foo × 100 · bar × 50", lifetime.Result)

	month, err := KicksMonth.Handler(context.Background(), parseCtx, &types.VariableData{Key: "top.kicks.month", Params: &params})
	require.NoError(t, err)
	require.Equal(t, "baz × 30", month.Result)

	week, err := KicksWeek.Handler(context.Background(), parseCtx, &types.VariableData{Key: "top.kicks.week", Params: &params})
	require.NoError(t, err)
	require.Equal(t, "qux × 10", week.Result)
}

func TestKickLeaderboardVariables_UnsupportedPlatform(t *testing.T) {
	parseCtx := &types.VariableParseContext{ParseContext: &types.ParseContext{Platform: platformentity.PlatformTwitch}}
	result, err := KicksLifetime.Handler(context.Background(), parseCtx, &types.VariableData{Key: "top.kicks.lifetime"})
	require.NoError(t, err)
	require.Equal(t, "not supported on this platform", result.Result)
}
