package stream

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/parser/internal/types"
	parserservices "github.com/twirapp/twir/apps/parser/internal/types/services"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorymodel "github.com/twirapp/twir/libs/repositories/channels_info_history/model"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

type fakeChannelsInfoHistoryRepo struct {
	history []channelsinfohistorymodel.ChannelInfoHistory
	err     error
	input   channelsinfohistory.GetManyInput
}

func (f *fakeChannelsInfoHistoryRepo) GetMany(ctx context.Context, input channelsinfohistory.GetManyInput) ([]channelsinfohistorymodel.ChannelInfoHistory, error) {
	f.input = input
	return f.history, f.err
}

func (f *fakeChannelsInfoHistoryRepo) Create(ctx context.Context, input channelsinfohistory.CreateInput) error {
	panic("unexpected call")
}

func TestCategoryTime_UsesCurrentStreamHistoryWindow(t *testing.T) {
	startedAt := time.Now().UTC().Add(-2 * time.Hour)
	categoryStartedAt := time.Now().UTC().Add(-45 * time.Minute)
	repo := &fakeChannelsInfoHistoryRepo{
		history: []channelsinfohistorymodel.ChannelInfoHistory{
			{ChannelID: "channel-1", Category: "Slots", CreatedAt: categoryStartedAt},
		},
	}

	parseCtx := &types.VariableParseContext{ParseContext: &types.ParseContext{
		Channel: &types.ParseContextChannel{ID: "channel-1"},
		ChannelStream: &streamsmodel.Stream{
			StartedAt: startedAt,
			GameName:  "Slots",
			ID:        "stream-1",
		},
		Services: &parserservices.Services{ChannelsInfoHistoryRepo: repo},
	}}

	result, err := CategoryTime.Handler(context.Background(), parseCtx, &types.VariableData{Key: "stream.category.time"})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotEmpty(t, result.Result)
	require.Equal(t, startedAt, repo.input.After)
}
