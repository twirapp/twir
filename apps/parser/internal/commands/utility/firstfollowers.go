package utility

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

type firstFollower struct {
	twitch.FirstFollower
	FollowedAt time.Time
}

var FirstFollowers = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "firstfollowers",
		Description: null.StringFrom("Shows channel first followers"),
		RolesIDS:    pq.StringArray{},
		Module:      "UTILITY",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		followers, err := twitch.GetFirstChannelFollowers(ctx, parseCtx.Channel.Name)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: err.Error(),
				Err:     err,
			}
		}

		filteredFollowers := make([]firstFollower, 0, len(followers))
		for idx, f := range followers {
			if f.Login == "" || f.DisplayName == "" {
				continue
			}

			filteredFollowers = append(
				filteredFollowers,
				firstFollower{
					FirstFollower: f,
				},
			)

			if idx == 10 {
				break
			}
		}

		var wg sync.WaitGroup
		for idx, f := range filteredFollowers {
			wg.Add(1)

			go func() {
				defer wg.Done()

				follow := parseCtx.Cacher.GetTwitchUserFollow(ctx, f.Id)
				if follow == nil {
					return
				}

				filteredFollowers[idx].FollowedAt = follow.Followed.Time
			}()
		}

		wg.Wait()

		formattedFollowers := make([]string, 0, len(filteredFollowers))
		for _, f := range filteredFollowers {
			time := helpers.Duration(
				f.FollowedAt,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Minutes: false,
						Seconds: true,
					},
				},
			)

			formattedFollowers = append(
				formattedFollowers,
				fmt.Sprintf(
					"%s (%s)",
					helpers.ResolveDisplayName(f.Login, f.DisplayName),
					strings.TrimSpace(time),
				),
			)
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				strings.Join(formattedFollowers, " · "),
			},
		}

		return result, nil
	},
}
