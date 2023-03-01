package user_songs_requested

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/samber/lo"
)

type sumResult struct {
	Sum int64 //or int ,or some else
}

var DurationVariable = types.Variable{
	Name:         "user.songs.requested.duration",
	Description:  lo.ToPtr("Duration of requested by user songs"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		sum := &sumResult{}
		err := db.
			Table("channels_requested_songs").
			Select("sum(duration) as sum").
			Where(`"channelId" = ? AND "orderedById" = ?`, ctx.ChannelId, ctx.SenderId).
			Scan(&sum).
			Error

		if err != nil {
			zap.S().Error(err)
			result.Result = "0"
			return result, nil
		}

		f := time.Duration(sum.Sum) * time.Millisecond
		result.Result = fmt.Sprintf("%.1fh", f.Hours())

		return result, nil
	},
}
