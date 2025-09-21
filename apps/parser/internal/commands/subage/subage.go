package subage

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
						i18n.GetCtx(
							ctx,
							locales.Translations.Commands.Subage.Errors.NotSubscriberOrHidden,
						),
					},
				}, nil
			}
			return nil, err
		}

		var resultMessage string

		if subAgeInfo.Cumulative == nil && subAgeInfo.Streak == nil {
			// User is not a subscriber
			resultMessage = i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Subage.Responses.NotSubscriber.SetVars(
					locales.KeysCommandsSubageResponsesNotSubscriberVars{
						User: userNameForCheck,
					},
				),
			)
		} else if subAgeInfo.Type == nil && subAgeInfo.Cumulative != nil && subAgeInfo.Cumulative.Months != nil {
			// User is not a subscriber but used to be
			resultMessage = i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Subage.Responses.NotSubscriberButWas.SetVars(
					locales.KeysCommandsSubageResponsesNotSubscriberButWasVars{
						User:   userNameForCheck,
						Months: *subAgeInfo.Cumulative.Months,
					},
				),
			)
		} else {
			// User has an active subscription
			var tier string
			if subAgeInfo.Type != nil {
				tier = subAgeInfo.Type.String()
			}

			var months int
			if subAgeInfo.Cumulative != nil && subAgeInfo.Cumulative.Months != nil {
				months = *subAgeInfo.Cumulative.Months
			}

			resultMessage = i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Subage.Responses.SubscriptionInfo.SetVars(
					locales.KeysCommandsSubageResponsesSubscriptionInfoVars{
						User:    userNameForCheck,
						Tier:    tier,
						Channel: parseCtx.Channel.Name,
						Months:  months,
					},
				),
			)

			// Add streak information if different from cumulative
			if subAgeInfo.Streak != nil &&
				subAgeInfo.Streak.Months != nil &&
				*subAgeInfo.Streak.Months != *subAgeInfo.Cumulative.Months {
				streakMessage := i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Subage.Responses.StreakInfo.SetVars(
						locales.KeysCommandsSubageResponsesStreakInfoVars{
							Months: *subAgeInfo.Streak.Months,
						},
					),
				)
				resultMessage += streakMessage
			}

			// Add time remaining if subscription has an end date
			if subAgeInfo.EndsAt != nil {
				duration := helpers.Duration(
					*subAgeInfo.EndsAt, &helpers.DurationOpts{
						UseUtc: true,
						Hide: helpers.DurationOptsHide{
							Seconds: true,
						},
					},
				)

				timeRemainingMessage := i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Subage.Responses.TimeRemaining.SetVars(
						locales.KeysCommandsSubageResponsesTimeRemainingVars{
							Duration: duration,
						},
					),
				)
				resultMessage += timeRemainingMessage
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{resultMessage},
		}, nil
	},
}
