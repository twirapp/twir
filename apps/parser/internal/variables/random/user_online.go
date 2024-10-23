package random

import (
	"context"
	"errors"
	"math/rand"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var OnlineUser = &types.Variable{
	Name:                "random.online.user",
	Description:         lo.ToPtr("Choose random online user"),
	CanBeUsedInRegistry: true,
	NotCachable:         true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		var onlineCount int64
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.UsersOnline{}).
			Where(`"channelId" = ? `, parseCtx.Channel.ID).
			Count(&onlineCount).Error
		if err != nil || onlineCount == 0 {
			return nil, errors.New("no users online")
		}

		randCount := rand.Intn(int(onlineCount)-0) + 0

		randomUser := &model.UsersOnline{}
		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? `, parseCtx.Channel.ID).
			Offset(randCount).
			First(randomUser).Error

		if err != nil || randomUser == nil {
			return nil, errors.New("cannot get online user")
		}

		result.Result = randomUser.UserName.String
		return result, nil
	},
}
