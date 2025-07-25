package dota

// import (
// 	"strconv"

// 	"github.com/guregu/null"
// 	"github.com/lib/pq"
// 	"github.com/samber/do"
// 	"github.com/twirapp/twir/apps/parser/internal/di"
// 	"go.uber.org/zap"
// 	"gorm.io/gorm"

// 	"github.com/twirapp/twir/apps/parser/internal/types"

// 	model "github.com/twirapp/twir/libs/gomodels"

// 	variables_cache "github.com/twirapp/twir/apps/parser/internal/variablescache"

// 	steamid "github.com/leighmacdonald/steamid/v2/steamid"
// 	"github.com/samber/lo"
// )

// var AddAccCommand = &types.DefaultCommand{
// 	ChannelsCommands: &model.ChannelsCommands{
// 		Name:        "dota addacc",
// 		Description: null.StringFrom("Add dota account for watching games"),
// 		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
// 		Module:      "DOTA",
// 		IsReply:     true,
// 	},
// 	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
// 		result := &types.CommandsHandlerResult{
// 			Result: make([]string, 0),
// 		}
// 		db := do.MustInvoke[gorm.DB](di.Provider)

// 		acc, err := strconv.ParseUint(*ctx.Text, 10, 64)
// 		if err != nil {
// 			result.Result = append(result.Result, WRONG_ACCOUNT_ID)
// 			return result
// 		}

// 		ok := lo.Try(func() error {
// 			n := steamid.SID32(acc)
// 			steamid.SID32ToSID(n)
// 			return nil
// 		})

// 		if !ok {
// 			result.Result = append(result.Result, WRONG_ACCOUNT_ID)
// 			return result
// 		}

// 		accId := steamid.SID32(acc)

// 		var count int64 = 0
// 		err = db.
// 			Table("channels_dota_accounts").
// 			Where(`"channelId" = ? AND "id" = ?`, ctx.ChannelId, strconv.Itoa(int(accId))).
// 			Count(&count).Error

// 		if err != nil {
// 			zap.S().Error(err)
// 			result.Result = append(result.Result, "Error happend on our side.")
// 			return result
// 		}

// 		if count != 0 {
// 			result.Result = append(result.Result, "Account already added.")
// 			return result
// 		}

// 		err = db.
// 			Create(&model.ChannelsDotaAccounts{
// 				ID:        strconv.Itoa(int(accId)),
// 				ChannelID: ctx.ChannelId,
// 			}).Error

// 		if err != nil {
// 			zap.S().Error(err)
// 			result.Result = append(
// 				result.Result,
// 				"Something went wrong on out side when inserting account into db.",
// 			)
// 			return result
// 		}

// 		result.Result = append(result.Result, "Account added.")
// 		return result
// 	},
// }
