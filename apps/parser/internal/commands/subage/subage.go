package subage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

var SubAge = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "subage",
		Description: null.StringFrom("Displays sub age of user or mentioned user."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		userNameForCheck := parseCtx.Sender.Name
		if len(parseCtx.Mentions) >= 1 {
			userNameForCheck = parseCtx.Mentions[0].UserName
		}

		subAgeInfo, err := parseCtx.Cacher.GetSubAgeInfo(ctx, parseCtx.Channel.Name, userNameForCheck)
		if err != nil {
			if errors.Is(err, twitch.ErrSubAgeEmpty) {
				return &types.CommandsHandlerResult{
					Result: []string{
						"not a subscriber or info hidden",
					},
				}, nil
			}
			return nil, err
		}

		var result strings.Builder

		result.WriteString(userNameForCheck)

		if subAgeInfo.Cumulative == nil && subAgeInfo.Streak == nil {
			result.WriteString(" is not a subscriber.")
		} else if subAgeInfo.Type == nil && subAgeInfo.Cumulative != nil && subAgeInfo.Cumulative.Months != nil {
			result.WriteString(" is not a subscriber, but used to be for ")
			result.WriteString(fmt.Sprint(*subAgeInfo.Cumulative.Months))
			result.WriteString(" months.")
		} else {
			result.WriteString(" has a ")
			if subAgeInfo.Type != nil {
				result.WriteString(subAgeInfo.Type.String())
				result.WriteString(" ")
			}
			result.WriteString("subscription to ")
			result.WriteString(parseCtx.Channel.Name)
			result.WriteString(" for a total ")
			if subAgeInfo.Cumulative != nil && subAgeInfo.Cumulative.Months != nil {
				result.WriteString(strconv.Itoa(*subAgeInfo.Cumulative.Months))
				result.WriteString(" months")
			}

			// if subAgeInfo.Cumulative != nil && subAgeInfo.Cumulative.ElapsedDays != nil {
			// 	result.WriteString(strconv.Itoa(*subAgeInfo.Cumulative.ElapsedDays))
			// 	result.WriteString(" days")
			// }

			if subAgeInfo.Streak != nil &&
				subAgeInfo.Streak.Months != nil &&
				*subAgeInfo.Streak.Months != *subAgeInfo.Cumulative.Months {
				result.WriteString(" , currently on a ")
				result.WriteString(strconv.Itoa(*subAgeInfo.Streak.Months))
				result.WriteString(" months streak")
			}

			if subAgeInfo.EndsAt != nil {
				duration := helpers.Duration(
					*subAgeInfo.EndsAt, &helpers.DurationOpts{
						UseUtc: true,
						Hide: helpers.DurationOptsHide{
							Seconds: true,
						},
					},
				)

				result.WriteString(", ")
				result.WriteString(fmt.Sprint(duration))
				result.WriteString(" remaining")
			}
		}

		//  L󠀀in󠀀ar󠀀yx󠀀 h󠀀as a tier 3 permanent subscription to L󠀀in󠀀ar󠀀yx for a total 23 months, currently on a 23 month streak.

		return &types.CommandsHandlerResult{
			Result: []string{result.String()},
		}, nil
	},
}
