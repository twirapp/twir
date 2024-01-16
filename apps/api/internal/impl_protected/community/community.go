package community

import (
	"context"
	"strings"

	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/community"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Community struct {
	*impl_deps.Deps
}

func (c *Community) CommunityResetStats(
	ctx context.Context,
	request *community.ResetStatsRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	if request.Field == community.ResetStatsRequest_Emotes {
		err := c.Db.WithContext(ctx).
			Where(`"channelId" = ?`, dashboardId).
			Delete(&model.ChannelEmoteUsage{}).Error
		if err != nil {
			return nil, err
		}

		return &emptypb.Empty{}, nil
	}

	field := strings.ToLower(request.Field.String())
	if request.Field == community.ResetStatsRequest_UsedChannelsPoints {
		field = "usedChannelPoints"
	}

	err := c.Db.WithContext(ctx).
		Model(&model.UsersStats{}).
		Where(`"channelId" = ?`, dashboardId).
		Update(field, 0).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
