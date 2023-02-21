package user_songs_requested

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
	"strconv"

	"github.com/samber/lo"
)

var CountVariable = types.Variable{
	Name:         "user.songs.requested.count",
	Description:  lo.ToPtr("How many songs user requested"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		var count int64
		err := db.
			Model(&model.RequestedSong{}).
			Where(`"channelId" =? AND "orderedById" = ?`, ctx.ChannelId, ctx.SenderId).
			Count(&count).
			Error

		if err != nil {
			result.Result = "0"
			return result, nil
		}

		result.Result = strconv.FormatInt(count, 10)

		return result, nil
	},
}
