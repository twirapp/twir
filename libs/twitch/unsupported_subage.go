package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/samber/lo"
)

var ErrSubAgeEmpty = fmt.Errorf("not subscribed or info hidden")

type ivrSubAgeResponse struct {
	User struct {
		Id          string `json:"id"`
		Login       string `json:"login"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	Channel struct {
		Id          string `json:"id"`
		Login       string `json:"login"`
		DisplayName string `json:"displayName"`
	} `json:"channel"`
	StatusHidden bool      `json:"statusHidden"`
	FollowedAt   time.Time `json:"followedAt"`
	Streak       *struct {
		ElapsedDays   int       `json:"elapsedDays"`
		DaysRemaining int       `json:"daysRemaining"`
		Months        int       `json:"months"`
		End           time.Time `json:"end"`
		Start         time.Time `json:"start"`
	} `json:"streak"`
	Cumulative *struct {
		ElapsedDays   int       `json:"elapsedDays"`
		DaysRemaining int       `json:"daysRemaining"`
		Months        int       `json:"months"`
		End           time.Time `json:"end"`
		Start         time.Time `json:"start"`
	} `json:"cumulative"`
	Meta *struct {
		Type     string     `json:"type"`
		Tier     string     `json:"tier"`
		EndsAt   *time.Time `json:"endsAt"`
		RenewsAt *time.Time `json:"renewsAt"`
		GiftMeta *struct {
			GiftDate time.Time `json:"giftDate"`
			Gifter   struct {
				Id          string `json:"id"`
				Login       string `json:"login"`
				DisplayName string `json:"displayName"`
			} `json:"gifter"`
		} `json:"giftMeta"`
	} `json:"meta"`
}

func GetUserSubAge(
	ctx context.Context,
	channelName, userName string,
) (*UserSubscribePayload, error) {
	apiUrl := fmt.Sprintf("https://api.ivr.fi/v2/twitch/subage/%s/%s", userName, channelName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get sub age: %d", res.StatusCode)
	}

	var requestData ivrSubAgeResponse
	if err := json.Unmarshal(body, &requestData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if requestData.StatusHidden {
		return nil, ErrSubAgeEmpty
	}

	if requestData.Streak == nil && requestData.Cumulative == nil {
		return nil, ErrSubAgeEmpty
	}

	payload := &UserSubscribePayload{
		Streak:     nil,
		Cumulative: nil,
		Type:       nil,
		Gifter:     nil,
		EndsAt:     nil,
		RenewsAt:   nil,
	}

	if requestData.Cumulative != nil {
		payload.Cumulative = &UserSubscribePayloadMeta{
			ElapsedDays:   &requestData.Cumulative.ElapsedDays,
			DaysRemaining: &requestData.Cumulative.DaysRemaining,
			Months:        &requestData.Cumulative.Months,
			End:           &requestData.Cumulative.End,
			Start:         &requestData.Cumulative.Start,
		}
	}

	if requestData.Streak != nil {
		payload.Streak = &UserSubscribePayloadMeta{
			ElapsedDays:   &requestData.Streak.ElapsedDays,
			DaysRemaining: &requestData.Streak.DaysRemaining,
			Months:        &requestData.Streak.Months,
			End:           &requestData.Streak.End,
			Start:         &requestData.Streak.Start,
		}
	}

	if requestData.Meta != nil {
		if requestData.Meta.GiftMeta != nil {
			payload.Gifter = &UserSubscribePayloadGifter{
				GifterID:       requestData.Meta.GiftMeta.Gifter.Id,
				GifterUsername: requestData.Meta.GiftMeta.Gifter.Login,
			}
		}
		switch requestData.Meta.Type {
		case "gift":
			payload.Type = lo.ToPtr(UserSubscribePayloadTypeGift)
		case "prime":
			payload.Type = lo.ToPtr(UserSubscribePayloadTypePrime)
		case "paid":
			payload.Type = lo.ToPtr(UserSubscribePayloadTypePaid)
		}

		if requestData.Meta.EndsAt != nil {
			payload.EndsAt = requestData.Meta.EndsAt
		}

		if requestData.Meta.RenewsAt != nil {
			payload.RenewsAt = requestData.Meta.RenewsAt
		}
	}

	return payload, nil
}

type UserSubscribePayloadMeta struct {
	ElapsedDays   *int
	DaysRemaining *int
	Months        *int
	End           *time.Time
	Start         *time.Time
}

type UserSubscribePayloadType string

func (c UserSubscribePayloadType) String() string {
	return string(c)
}

const (
	UserSubscribePayloadTypeGift  UserSubscribePayloadType = "gift"
	UserSubscribePayloadTypePrime UserSubscribePayloadType = "prime"
	UserSubscribePayloadTypePaid  UserSubscribePayloadType = "paid"
)

type UserSubscribePayloadGifter struct {
	GifterID       string
	GifterUsername string
}

type UserSubscribePayload struct {
	Streak     *UserSubscribePayloadMeta
	Cumulative *UserSubscribePayloadMeta

	Type   *UserSubscribePayloadType
	Gifter *UserSubscribePayloadGifter

	EndsAt   *time.Time
	RenewsAt *time.Time
}
