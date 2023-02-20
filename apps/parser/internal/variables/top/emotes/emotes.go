package emotes

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Emote struct {
	Emote string
	Count int

	UserID string
}

var Variable = types.Variable{
	Name:        "top.emotes",
	Description: lo.ToPtr("Top used emotes"),
	Example:     lo.ToPtr("top.emotes|10"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		db := do.MustInvoke[gorm.DB](di.Provider)
		result := &types.VariableHandlerResult{}

		limit := 10
		if data.Params != nil {
			newLimit, err := strconv.Atoi(*data.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		emotes := []Emote{}
		err := db.
			Raw(`SELECT emote, COUNT(*) 
				FROM channels_emotes_usages
				WHERE "channelId" = ?
				Group By emote 
				Order By COUNT(*) 
				DESC LIMIT ?
			`, ctx.ChannelId, limit).
			Scan(&emotes).
			Error

		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(emotes, func(e Emote, _ int) string {
			return fmt.Sprintf(
				"%s × %v",
				e.Emote,
				e.Count,
			)
		})

		result.Result = strings.Join(mappedTop, " ")
		return result, nil
	},
}
