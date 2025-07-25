package games

import (
	"context"
	"fmt"
	"math"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

var DuelStats = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "duel stats",
		Description: null.StringFrom("User stats for duels"),
		Module:      "GAMES",
		IsReply:     false,
		Visible:     true,
		Enabled:     false,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var duelsCount int64
		var winsCount int64
		var losesCount int64

		if err := parseCtx.Services.Gorm.WithContext(ctx).Model(&model.ChannelDuel{}).Where(
			`channel_id = ? and finished_at is not null and (sender_id = ? OR target_id = ?)`,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
			parseCtx.Sender.ID,
		).Count(&duelsCount).Error; err != nil {
			return nil, err
		}

		if err := parseCtx.Services.Gorm.WithContext(ctx).Model(&model.ChannelDuel{}).Where(
			`channel_id = ? and finished_at is not null and (sender_id = ? OR target_id = ?) and loser_id != ?`,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
			parseCtx.Sender.ID,
			parseCtx.Sender.ID,
		).Count(&winsCount).Error; err != nil {
			return nil, err
		}

		losesCount = duelsCount - winsCount

		winRate := float64(winsCount) / float64(duelsCount) * 100
		if math.IsNaN(winRate) {
			winRate = 0
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"You have shoot %d times · %d W – %d L (%.0f%% WR)",
					duelsCount,
					winsCount,
					losesCount,
					winRate,
				),
			},
		}, nil
	},
}
