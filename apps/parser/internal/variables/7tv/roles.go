package seventv

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Roles = &types.Variable{
	Name:         "7tv.roles",
	Description:  lo.ToPtr("Roles of user on 7tv"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile, err := parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			result.Result = fmt.Sprintf("[Twir err] Failed to get 7tv profile: %s", err)
			return &result, nil
		}

		roles := make([]string, 0, len(profile.Roles))
		for _, role := range profile.Roles {
			if role.Name == "Default" {
				continue
			}

			roles = append(roles, role.Name)
		}

		if len(roles) == 0 {
			result.Result = "No roles"
		} else {
			result.Result = strings.Join(roles, ", ")
		}

		return &result, nil
	},
}
