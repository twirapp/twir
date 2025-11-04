package top

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/cache/twitch"
)

type emotesUsersRow struct {
	UserID string `gorm:"column:userId" json:"userId" db:"userId"`
	Emotes int    `gorm:"column:emotes" json:"emotes" db:"emotes"`
}

var EmotesUsers = &types.Variable{
	Name:                "top.emotes.users",
	Description:         lo.ToPtr("Top users by emotes. Prints counts of used emotes"),
	Example:             lo.ToPtr("top.emotes.users|10"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		var page = 1

		if parseCtx.Text != nil {
			p, err := strconv.Atoi(*parseCtx.Text)
			if err == nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		query := `
SELECT "userId", "emotes"
FROM users_stats
WHERE "channelId" = ? AND emotes > 0
ORDER BY emotes DESC
LIMIT ?
OFFSET ?
`

		usages := make([]emotesUsersRow, 0, limit)

		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Raw(
				query, parseCtx.Channel.ID, limit, limit*(page-1),
			).
			Scan(&usages).
			Error
		if err != nil {
			return nil, err
		}

		usersIds := make([]string, 0, len(usages))
		for _, usage := range usages {
			usersIds = append(usersIds, usage.UserID)
		}

		twitchUsers, err := parseCtx.Services.CacheTwitchClient.GetUsersByIds(ctx, usersIds)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return nil, err
		}

		mappedTop := make([]string, 0, len(usages))

		for _, usage := range usages {
			user, ok := lo.Find(
				twitchUsers, func(item twitch.TwitchUser) bool {
					return item.ID == usage.UserID
				},
			)

			if !ok {
				continue
			}

			mappedTop = append(
				mappedTop,
				fmt.Sprintf(
					"%s × %v",
					user.Login,
					usage.Emotes,
				),
			)
		}

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
