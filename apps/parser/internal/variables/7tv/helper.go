package seventv

import (
	"context"

	"github.com/twirapp/twir/apps/parser/internal/types"
	seventvintegrationapi "github.com/twirapp/twir/libs/integrations/seventv/api"
)

func getProfile(ctx context.Context, parseCtx *types.VariableParseContext) (*seventvintegrationapi.TwirSeventvUser, error) {
	if parseCtx.Platform == "kick" {
		return parseCtx.Cacher.GetSeventvProfileGetKickId(ctx, parseCtx.Channel.ID)
	}

	return parseCtx.Cacher.GetSeventvProfileGetTwitchId(ctx, parseCtx.Channel.ID)
}
