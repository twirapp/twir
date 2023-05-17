package channel_title

import (
	"context"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

var History = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "title history",
		Description: null.StringFrom("Print history of titles."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		limit := 5

		if parseCtx.Text != nil && *parseCtx.Text != "" {
			l, err := strconv.Atoi(*parseCtx.Text)
			if err == nil {
				limit = l
			}
		}

		if limit > 20 {
			limit = 5
		}

		var histories []*model.ChannelInfoHistory
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Raw(`SELECT * FROM (
				SELECT DISTINCT ON (title) * FROM "channels_info_history" 
				                             WHERE "channelId" = ?
				                             ORDER BY "title", "createdAt" 
				                             DESC
				) subquery ORDER BY "createdAt" DESC LIMIT ?`, parseCtx.Channel.ID, limit).
			Find(&histories).
			Error

		if err != nil {
			result.Result = append(result.Result, "internal error")
			return result
		}

		titles := lo.Map(histories, func(item *model.ChannelInfoHistory, _ int) string {
			return item.Title
		})

		result.Result = append(result.Result, strings.Join(titles, " ║ "))
		return result
	},
}
