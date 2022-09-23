package dota

import (
	"fmt"
	"strconv"
	"tsuwari/parser/internal/types"

	model "tsuwari/models"

	variables_cache "tsuwari/parser/internal/variablescache"

	steamid "github.com/leighmacdonald/steamid/v2/steamid"
	"github.com/samber/lo"
)

var DelAccCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "dota delacc",
		Description: lo.ToPtr("Delete dota account "),
		Permission:  "BROADCASTER",
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		acc, err := strconv.ParseUint(*ctx.Text, 10, 64)
		if err != nil {
			return []string{WRONG_ACCOUNT_ID}
		}

		ok := lo.Try(func() error {
			n := steamid.SID32(acc)
			steamid.SID32ToSID(n)
			return nil
		})

		if !ok {
			return []string{WRONG_ACCOUNT_ID}
		}

		accId := steamid.SID32(acc)

		var count int64 = 0
		err = ctx.Services.Db.
			Table("channels_dota_accounts").
			Where(`"channelId" = ? AND "id" = ?`, ctx.ChannelId, strconv.Itoa(int(accId))).
			Count(&count).Error

		if err != nil {
			fmt.Println(err)
			return []string{"Error happend on our side."}
		}

		if count == 0 {
			return []string{"Account not added."}
		}

		err = ctx.Services.Db.
			Delete(&model.ChannelsDotaAccounts{
				ID:        strconv.Itoa(int(accId)),
				ChannelID: ctx.ChannelId,
			}).Error

		if err != nil {
			fmt.Println(err)
			return []string{"Something went wrong on out side when inserting account into db."}
		}

		return []string{"Account removed."}
	},
}
