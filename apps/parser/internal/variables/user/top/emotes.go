package user_top

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"gorm.io/gorm"
	"strings"
)

type Emote struct {
	Emote string
	Count int
}

var TopEmotesVariable = types.Variable{
	Name:         "user.top.emotes",
	Description:  lo.ToPtr("User top used emotes"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		emotes := []Emote{}
		err := db.
			Raw(`SELECT emote, COUNT(*) 
				FROM channels_emotes_usages
				WHERE "channelId" = ? AND "userId" = ?
				Group By emote 
				Order By COUNT(*) 
				DESC LIMIT 20
			`, ctx.ChannelId, ctx.SenderId).
			Scan(&emotes).
			Error

		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(emotes, func(e Emote, idx int) string {
			return fmt.Sprintf(
				"%v. %s — %v",
				idx+1,
				e.Emote,
				e.Count,
			)
		})

		result.Result = strings.Join(mappedTop, ", ")
		return result, nil
	},
}
