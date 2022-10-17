package randomonlineuser

import (
	"errors"
	"math/rand"
	"time"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

type OnlineUser struct {
	UserName  string
	UserId    string
	ChannelId string
}

var Variable = types.Variable{
	Name:        "random.online.user",
	Description: lo.ToPtr("Choose random online user"),
	Example:     lo.ToPtr("random.online.user"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		onlineCount := int64(0)
		err := ctx.Services.Db.
			Model(&model.UsersOnline{}).
			Where(`"channelId" = ? `, ctx.ChannelId).
			Count(&onlineCount).Error
		if err != nil || onlineCount == 0 {
			return nil, errors.New("no users online")
		}

		rand.Seed(time.Now().Unix())
		randCount := rand.Intn(int(onlineCount)-0) + 0

		randomUser := &model.UsersOnline{}
		err = ctx.Services.Db.
			Model(&model.UsersOnline{}).
			Offset(randCount).
			First(randomUser).Error

		if err != nil || randomUser == nil {
			return nil, errors.New("cannot get online user")
		}

		result.Result = randomUser.UserName.String
		return result, nil
	},
}
