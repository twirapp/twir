package integrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_seventv"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"github.com/twitchtv/twirp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (c *Integrations) IntegrationsSevenTvGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_seventv.GetDataResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var defaultBot model.Bots
	if err := c.Db.
		WithContext(ctx).
		Where(`type = ?`, "DEFAULT").
		First(&defaultBot).
		Error; err != nil {
		return nil, err
	}

	var sevenTvSettings model.ChannelsIntegrationsSettingsSeventv
	err = c.Db.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		First(&sevenTvSettings).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		sevenTvSettings.ID = uuid.New()
		sevenTvSettings.ChannelID = dashboardId
		sevenTvSettings.RewardIdForAddEmote = null.StringFromPtr(nil)
		sevenTvSettings.RewardIdForRemoveEmote = null.StringFromPtr(nil)
		sevenTvSettings.DeleteEmotesOnlyAddedByApp = false
		sevenTvSettings.AddedEmotes = pq.StringArray{}

		if err := c.Db.
			WithContext(ctx).
			Save(&sevenTvSettings).
			Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	var botSevenTvResponse seventv.ProfileResponse
	var userSevenTvResponse seventv.ProfileResponse

	wg, wgCtx := errgroup.WithContext(ctx)
	wg.Go(
		func() error {
			resp, err := seventv.GetProfile(wgCtx, defaultBot.ID)
			if err != nil {
				c.Logger.Error("failed to get 7tv bot data", "err", err)
				return err
			}

			botSevenTvResponse = resp
			return nil
		},
	)

	wg.Go(
		func() error {
			resp, err := seventv.GetProfile(wgCtx, dashboardId)
			if err != nil {
				c.Logger.Error("failed to get 7tv channel data", "err", err)
				return err
			}

			userSevenTvResponse = resp
			return nil
		},
	)

	if err := wg.Wait(); err != nil {
		if errors.Is(err, errSevenTvProfileNotFound) {
			return nil, twirp.NewError(twirp.NotFound, "profile_not_found")
		}

		return nil, err
	}

	var isBotEditor bool
	if botSevenTvResponse.Id == userSevenTvResponse.Id {
		isBotEditor = true
	} else {
		for _, editor := range userSevenTvResponse.User.Editors {
			if editor.Id == botSevenTvResponse.User.Id {
				isBotEditor = true
				break
			}
		}
	}

	resp := &integrations_seventv.GetDataResponse{
		IsEditor: isBotEditor,
		BotSeventvProfile: &integrations_seventv.SevenTvProfile{
			Id:          botSevenTvResponse.User.Id,
			Username:    botSevenTvResponse.User.Username,
			DisplayName: botSevenTvResponse.User.DisplayName,
		},
		UserSeventvProfile: &integrations_seventv.SevenTvProfile{
			Id:          userSevenTvResponse.User.Id,
			Username:    userSevenTvResponse.User.Username,
			DisplayName: userSevenTvResponse.User.DisplayName,
		},
		RewardIdForAddEmote:        sevenTvSettings.RewardIdForAddEmote.Ptr(),
		RewardIdForRemoveEmote:     sevenTvSettings.RewardIdForRemoveEmote.Ptr(),
		DeleteEmotesOnlyAddedByApp: sevenTvSettings.DeleteEmotesOnlyAddedByApp,
	}

	if userSevenTvResponse.EmoteSet != nil {
		resp.EmoteSetId = &userSevenTvResponse.EmoteSet.Id
	}

	return resp, nil
}

func (c *Integrations) IntegrationsSevenTvUpdate(
	ctx context.Context,
	req *integrations_seventv.UpdateDataRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var sevenTvSettings model.ChannelsIntegrationsSettingsSeventv
	if err := c.Db.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		First(&sevenTvSettings).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sevenTvSettings.ID = uuid.New()
		} else {
			return nil, err
		}
	}

	sevenTvSettings.ChannelID = dashboardId
	sevenTvSettings.RewardIdForAddEmote = null.StringFromPtr(req.RewardIdForAddEmote)       //nolint:protogetter
	sevenTvSettings.RewardIdForRemoveEmote = null.StringFromPtr(req.RewardIdForRemoveEmote) //nolint:protogetter
	sevenTvSettings.DeleteEmotesOnlyAddedByApp = req.GetDeleteEmotesOnlyAddedByApp()

	if err := c.Db.
		WithContext(ctx).
		Save(&sevenTvSettings).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

var errSevenTvProfileNotFound = fmt.Errorf("7tv profile not found")
