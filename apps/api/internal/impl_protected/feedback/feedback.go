package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/api/messages/feedback"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Feedback struct {
	*impl_deps.Deps
}

func getRedisKey(ip string) string {
	return "feedback-limit:" + ip
}

func (c *Feedback) LeaveFeedback(
	ctx context.Context,
	req *feedback.LeaveFeedbackRequest,
) (*emptypb.Empty, error) {
	message := req.GetMessage()

	if len(message) == 0 || len(message) > 1000 {
		return nil, twirp.InvalidArgumentError(
			"message",
			"message cannot be empty and greater than 1000 symbols",
		)
	}

	headers, err := helpers.GetHeadersFromCtx(ctx)
	if err != nil {
		return nil, twirp.Unavailable.Error("cannot get headers of request")
	}

	var ip string

	headerRealIp := headers.Get("x-real-ip")
	if headerRealIp != "" {
		ip = headerRealIp
	} else {
		ip = headers.Get("x-forwarded-for")
	}
	if ip == "" {
		return nil, twirp.Internal.Error("real ip not provided")
	}

	isRateLimited, err := c.Redis.Exists(ctx, getRedisKey(ip)).Result()
	if err != nil {
		return nil, twirp.Internal.Error("cannot get rate limit value")
	}

	if isRateLimited > 0 {
		d, err := c.Redis.TTL(ctx, getRedisKey(ip)).Result()
		if err != nil {
			return nil, twirp.Internal.Error("cannot get ttl of rate limit")
		}

		return nil, twirp.ResourceExhausted.
			Error("rate limited").
			WithMeta("retryable", "true").
			WithMeta("retry_after", d.String())
	}

	user, err := helpers.GetUserModelFromCtx(ctx)
	if err != nil {
		return nil, twirp.Internal.Error("cannot find user in request")
	}

	if c.Config.DiscordFeedbackUrl == "" {
		return nil, twirp.Unavailable.Error("discord not configured on backend")
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, twirp.Internal.Error("cannot create twitch client")
	}

	helixUserRequest, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{user.ID},
		},
	)
	if err != nil {
		return nil, twirp.Internal.Error("cannot request user from twitch")
	}
	if helixUserRequest.ErrorMessage != "" {
		return nil, twirp.Internal.Error("cannot request user from twitch")
	}

	if len(helixUserRequest.Data.Users) == 0 {
		return nil, twirp.NotFound.Error("cannot find user on twitch")
	}

	helixUser := helixUserRequest.Data.Users[0]
	err = c.sendEmbed(
		sendEmbedOpts{
			authorName:   helixUser.Login,
			authorAvatar: helixUser.ProfileImageURL,
			authorID:     helixUser.ID,
			message:      message,
		},
	)
	if err != nil {
		return nil, twirp.Internal.Error(err.Error())
	}

	c.Redis.Set(ctx, getRedisKey(ip), "", 15*time.Minute)

	return &emptypb.Empty{}, nil
}

type embedAuthor struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type embedThumbNail struct {
	Url string `json:"url"`
}

type discordEmbed struct {
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Author      embedAuthor    `json:"author"`
	ThumbNail   embedThumbNail `json:"thumbnail"`
}

type discordRequest struct {
	Embeds  []discordEmbed `json:"embeds"`
	Content string         `json:"content"`
}

type sendEmbedOpts struct {
	authorName   string
	authorAvatar string
	authorID     string
	message      string
}

func (c *Feedback) sendEmbed(
	opts sendEmbedOpts,
) error {
	embed := discordEmbed{
		Type:        "rich",
		Title:       "New feedback",
		Description: opts.message,
		Color:       0x00FFFF,
		Author: embedAuthor{
			Name:    fmt.Sprintf("%s#%s", opts.authorName, opts.authorID),
			IconUrl: opts.authorAvatar,
		},
		ThumbNail: embedThumbNail{
			Url: opts.authorAvatar,
		},
	}

	requestBytes, err := json.Marshal(
		&discordRequest{
			Embeds:  []discordEmbed{embed},
			Content: "New feedback",
		},
	)
	if err != nil {
		return fmt.Errorf("wrong request body: %w", err)
	}

	resp, err := http.Post(
		c.Config.DiscordFeedbackUrl,
		"application/json",
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return fmt.Errorf("cannot send request to discord: %w", err)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("cannot send request to discord with status %v", resp.StatusCode)
	}

	return nil
}
