package user_emotes

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"gorm.io/gorm"
)

var Variable = types.Variable{
	Name:         "user.emotes",
	Description:  lo.ToPtr("User used emotes count"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		var count int64
		err := db.Where(`"channelId" = ? AND "userId" = ?`, ctx.ChannelId, ctx.SenderId).Count(&count).Error

		if err != nil {
			fmt.Println(err)
			result.Result = "error"
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
