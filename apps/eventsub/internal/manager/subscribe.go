package manager

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go/v4"
	eventsub_framework "github.com/twirapp/twitch-eventsub-framework"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Manager) SubscribeWithLimits(
	ctx context.Context,
	srq *eventsub_framework.SubRequest,
) (
	*esb.RequestStatus,
	error,
) {
	data, err := retry.DoWithData(
		func() (*esb.RequestStatus, error) {
			data, subscribeErr := c.Subscribe(ctx, srq)

			return data, subscribeErr
		},
		retry.Attempts(0),
		retry.Delay(1*time.Second),
		retry.RetryIf(
			func(err error) bool {
				var e *eventsub_framework.TwitchError
				if errors.As(err, &e) && e.Status != 409 {
					if e.Status == 429 {
						return true
					}
				}

				return false
			},
		),
	)

	return data, err
}
