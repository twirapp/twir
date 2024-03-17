package integrations

import (
	"context"

	"github.com/NovikovRoman/pubg"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_pubg"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Integrations) IntegrationsPubgGet(
	ctx context.Context, _ *emptypb.Empty,
) (*integrations_pubg.GetDataResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServicePubg,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_pubg.GetDataResponse{
		Nickname: integration.Data.UserName,
	}, nil
}

func (c *Integrations) IntegrationsPubgPut(
	ctx context.Context,
	req *integrations_pubg.PostDataRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServicePubg,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	players, err := c.PubgClient.GetPlayerByNickname(ctx, req.GetNickname())
	if err != nil {
		if _, ok := err.(*pubg.ErrNotFound); ok {
			return nil, twirp.NewError(twirp.NotFound, "player not found")
		}

		return nil, err
	}

	player := players.Data[0]

	integration.Enabled = true
	integration.Data.UserName = lo.ToPtr(req.GetNickname())
	integration.Data.UserId = lo.ToPtr(player.ID)
	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
