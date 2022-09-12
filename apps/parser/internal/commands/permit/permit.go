package permit

import (
	"fmt"
	"strconv"
	"strings"
	model "tsuwari/parser/internal/models"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/nicklaw5/helix"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "permit",
		Description: lo.ToPtr("Permits user."),
		Permission:  "MODERATOR",
		Visible:     true,
		Module:      lo.ToPtr("CHANNEL"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		count := 1
		params := strings.Split(*ctx.Text, " ")

		paramsLen := len(params)
		if paramsLen < 1 {
			return []string{"you have type user name to permit."}
		}

		if paramsLen == 2 {
			newCount, err := strconv.Atoi(params[1])
			if err == nil {
				count = newCount
			}
		}

		if count > 100 {
			return []string{"cannot create more then 100 permits."}
		}

		target, err := ctx.Services.Twitch.Client.GetUsers(&helix.UsersParams{
			Logins: []string{params[0]},
		})

		if err != nil || target.StatusCode != 200 || len(target.Data.Users) == 0 {
			return []string{"user not found."}
		}

		ctx.Services.Db.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < count; i++ {
				permit := model.ChannelsPermits{
					ID:        uuid.NewV4().String(),
					ChannelID: ctx.ChannelId,
					UserID:    target.Data.Users[0].ID,
				}
				err := tx.Create(&permit).Error
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
			return nil
		})

		return []string{fmt.Sprintf("✅ added %v permits to %s", count, target.Data.Users[0].DisplayName)}
	},
}
