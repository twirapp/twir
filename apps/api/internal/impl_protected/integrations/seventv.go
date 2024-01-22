package integrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_seventv"
	"github.com/twitchtv/twirp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
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
		Where("id = ? AND type = ?", dashboardId, "DEFAULT").
		First(&defaultBot).
		Error; err != nil {
		return nil, err
	}

	var sevenTvSettings model.ChannelsIntegrationsSettingsSeventv
	if err := c.Db.
		WithContext(ctx).
		Where("channel_id = ?", dashboardId).
		Find(&sevenTvSettings).
		Error; err != nil {
		return nil, err
	}

	var botSevenTvResponse sevenTvResponse
	var userSevenTvResponse sevenTvResponse

	wg, wgCtx := errgroup.WithContext(ctx)
	wg.Go(
		func() error {
			resp, err := c.getSevenTvDataById(wgCtx, defaultBot.ID)
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
			resp, err := c.getSevenTvDataById(wgCtx, dashboardId)
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
			if editor.Id == defaultBot.ID {
				isBotEditor = true
				break
			}
		}
	}

	return &integrations_seventv.GetDataResponse{
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
		RewardIdForAddEmote:    sevenTvSettings.RewardIdForAddEmote.Ptr(),
		RewardIdForRemoveEmote: sevenTvSettings.RewardIdForRemoveEmote.Ptr(),
	}, nil
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
		Find(&sevenTvSettings).
		Error; err != nil {
		return nil, err
	}

	if sevenTvSettings.ID.String() == "" {
		sevenTvSettings.ID = uuid.New()
	}

	sevenTvSettings.ChannelID = dashboardId
	sevenTvSettings.RewardIdForAddEmote = null.StringFromPtr(req.RewardIdForAddEmote)
	sevenTvSettings.RewardIdForRemoveEmote = null.StringFromPtr(req.RewardIdForRemoveEmote)

	if err := c.Db.
		WithContext(ctx).
		Save(&sevenTvSettings).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

var errSevenTvProfileNotFound = fmt.Errorf("7tv profile not found")

func (c *Integrations) getSevenTvDataById(ctx context.Context, userId string) (
	sevenTvResponse,
	error,
) {
	var response sevenTvResponse
	resp, err := req.
		SetContext(ctx).
		SetSuccessResult(&response).
		Get("https://7tv.io/v3/users/twitch/" + userId)
	if err != nil {
		return response, err
	}
	if !resp.IsSuccessState() {
		if resp.StatusCode == 404 {
			return response, errSevenTvProfileNotFound
		}
		return response, fmt.Errorf("failed to get 7tv data: %s", resp.String())
	}

	return response, nil
}

type sevenTvResponse struct {
	Id            string      `json:"id"`
	Platform      string      `json:"platform"`
	Username      string      `json:"username"`
	DisplayName   string      `json:"display_name"`
	LinkedAt      int64       `json:"linked_at"`
	EmoteCapacity int         `json:"emote_capacity"`
	EmoteSetId    interface{} `json:"emote_set_id"`
	EmoteSet      struct {
	} `json:"emote_set"`
	User struct {
		Id          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		CreatedAt   int64  `json:"created_at"`
		AvatarUrl   string `json:"avatar_url"`
		Style       struct {
		} `json:"style"`
		Editors []struct {
			Id          string `json:"id"`
			Permissions int    `json:"permissions"`
			Visible     bool   `json:"visible"`
			AddedAt     int64  `json:"added_at"`
		} `json:"editors"`
		Roles       []string `json:"roles"`
		Connections []struct {
			Id            string      `json:"id"`
			Platform      string      `json:"platform"`
			Username      string      `json:"username"`
			DisplayName   string      `json:"display_name"`
			LinkedAt      int64       `json:"linked_at"`
			EmoteCapacity int         `json:"emote_capacity"`
			EmoteSetId    interface{} `json:"emote_set_id"`
			EmoteSet      struct {
				Id         string        `json:"id"`
				Name       string        `json:"name"`
				Flags      int           `json:"flags"`
				Tags       []interface{} `json:"tags"`
				Immutable  bool          `json:"immutable"`
				Privileged bool          `json:"privileged"`
				Capacity   int           `json:"capacity"`
				Owner      interface{}   `json:"owner"`
			} `json:"emote_set"`
		} `json:"connections"`
	} `json:"user"`
}
