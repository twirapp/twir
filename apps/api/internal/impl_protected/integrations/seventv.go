package integrations

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_seventv"
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
		BotUsername:        botSevenTvResponse.User.Username,
		BotUserDisplayName: botSevenTvResponse.User.DisplayName,
		IsEditor:           isBotEditor,
	}, nil
}

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
