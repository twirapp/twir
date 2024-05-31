package manager

import (
	"context"
	"errors"
	"strings"
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
			return c.Subscribe(ctx, srq)
		},
		retry.Attempts(0),
		retry.Delay(1*time.Second),
		retry.RetryIf(
			func(err error) bool {
				var e *eventsub_framework.TwitchError
				if errors.As(err, &e) && e.Status == 429 {
					return true
				}

				if errors.Is(err, context.DeadlineExceeded) {
					return true
				}

				if strings.Contains(err.Error(), "context deadline exceeded") {
					return true
				}

				return false
			},
		),
	)

	return data, err
}
