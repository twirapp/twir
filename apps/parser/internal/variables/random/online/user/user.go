package randomonlineuser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
	"tsuwari/parser/internal/types"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

type OnlineUser struct {
	UserName string
	UserId string
	ChannelId string
}

var Variable = types.Variable{
	Name:        "random.online.user",
	Description: lo.ToPtr("Choose random online user"),
	Example:     lo.ToPtr("random.online.user"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		
		onlineUsersKeys, err := ctx.Services.Redis.
			Keys(context.TODO(), fmt.Sprintf("onlineUsers:%s:*", ctx.ChannelId)).
			Result()

		if err != nil {
			log.Fatal(err)
			result.Result = "cannot fetch online users"
			return result, nil
		}

		onlineUsers, err := ctx.Services.Redis.
			MGet(context.TODO(), onlineUsersKeys...).
			Result()

		if err != nil {
			log.Fatal(err)
			result.Result = "cannot fetch online users"
			return result, nil
		}

		users := make([]OnlineUser, len(onlineUsersKeys))
		for i, user := range onlineUsers {
			parsedUser := OnlineUser{}
	
			err := json.Unmarshal([]byte(user.(string)), &parsedUser)
	
			if err == nil {
				users[i] = parsedUser
			}
		}
		rand.Seed(time.Now().Unix())
		user := users[rand.Int() % len(users)]
		result.Result = user.UserName
		return result, nil
	},
}
