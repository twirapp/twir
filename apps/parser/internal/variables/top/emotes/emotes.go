package emotes

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

var Variable = types.Variable{
	Name:        "top.emotes",
	Description: lo.ToPtr("Top emotes"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		emotes := []Emote{}
		err := db.
			Raw(`SELECT emote, COUNT(*) 
				FROM channels_emotes_usages
				WHERE "channelId" = ?
				Group By emote 
				Order By COUNT(*) 
				DESC LIMIT 30
			`, ctx.ChannelId).
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
