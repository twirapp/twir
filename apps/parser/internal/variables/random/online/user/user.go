package randomonlineuser

import (
	"errors"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

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
		db := do.MustInvoke[gorm.DB](di.Provider)

		onlineCount := int64(0)
		err := db.
			Model(&model.UsersOnline{}).
			Where(`"channelId" = ? `, ctx.ChannelId).
			Count(&onlineCount).Error
		if err != nil || onlineCount == 0 {
			return nil, errors.New("no users online")
		}

		rand.Seed(time.Now().Unix())
		randCount := rand.Intn(int(onlineCount)-0) + 0

		randomUser := &model.UsersOnline{}
		err = db.
			Where(`"channelId" = ? `, ctx.ChannelId).
			Offset(randCount).
			First(randomUser).Error

		if err != nil || randomUser == nil {
			return nil, errors.New("cannot get online user")
		}

		result.Result = randomUser.UserName.String
		return result, nil
	},
}
